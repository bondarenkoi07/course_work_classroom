package Service

import (
	"awesomeProject1/ClassroomDAO"
	ic "awesomeProject1/IntelligentClassroom"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"regexp"
	"sync"
)

// Struct Service implements standard MVC view for app.
// Methods of Service manipulates data that come form client,
// validate it and send a response
// deviceDAO responsible for access to classroom devices work logic
// upgrader responsible for data transporting in JSON-format without page reloading
type Service struct {
	distributor      chan *ic.ClassroomController
	clients          map[*websocket.Conn]bool
	deviceDAO        ClassroomDAO.DevicesDAO
	upgrader         websocket.Upgrader
	errorDistributor chan Error
	m                sync.Mutex
}

type Error struct {
	Error  error           `json:"error"`
	Client *websocket.Conn `json:"-"`
}

//Constructor of Service
// NewService() initiate structure's params by default values
// returns object of Service
func NewService() Service {
	var service = Service{}
	service.deviceDAO.Init()
	service.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	service.upgrader.CheckOrigin = func(r *http.Request) bool {
		validate, err := regexp.MatchString("^(192.168.1.196):([0-9]*)", r.RemoteAddr)
		if err != nil {
			log.Print("something goes wrong " + r.RemoteAddr)
		}
		if validate {
			return true
		} else {
			return false
		}
	}
	service.distributor = make(chan *ic.ClassroomController)
	service.clients = make(map[*websocket.Conn]bool)
	service.errorDistributor = make(chan Error)
	return service
}

//return view of start page via static files
func (s *Service) Index() http.Handler {
	return http.FileServer(http.Dir("static"))
}

//Action receives WevSocket data
//And delegate it to deviceDao
//If connection with client is refused
//or error happened while data transporting
//then Action abort it
func (s *Service) Action(w http.ResponseWriter, r *http.Request) {
	conn, err := (*s).upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WTF3")
		log.Println(err)
		return
	}
	(*s).clients[conn] = true
	defer conn.Close()

	err = conn.WriteJSON((*s).deviceDAO.GetRepo())
	if err != nil {
		log.Println("WTF1")
		(*s).clients[conn] = false
		log.Println(err)
		return
	}

	for {
		var msg = make(map[string]string)
		err := conn.ReadJSON(&msg)
		//log.Print("Got message: ", msg, "\n")
		if err != nil {
			log.Println("WTF2")
			(*s).clients[conn] = false
			log.Println(err)
			return
		}
		actionValue, statusAction := msg["action"]
		deviceValue, statusDevice := msg["device"]
		if statusAction && statusDevice {
			switch actionValue {
			case "activate":
				err := (*s).deviceDAO.ActivateDevice(deviceValue)
				if err != nil {
					(*s).errorDistributor <- Error{Error: err, Client: conn}
				}
				break
			case "deactivate":
				err := (*s).deviceDAO.DeactivateDevice(deviceValue)
				if err != nil {
					(*s).errorDistributor <- Error{Error: err, Client: conn}
				}
				break
			default:
				(*s).errorDistributor <- Error{Error: err, Client: conn}
				return
			}
		} else {
			(*s).errorDistributor <- Error{Error: err, Client: conn}
			return
		}
		(*s).distributor <- (*s).deviceDAO.GetRepo()
	}

}

func (s *Service) ClassroomControllerDistribution() {
	for {
		message := <-(*s).distributor
		log.Println(message)
		for client := range (*s).clients {
			var err error
			if (*s).clients[client] {
				s.m.Lock()
				err = client.WriteJSON(message)
				s.m.Unlock()
			}
			if err != nil {
				client.Close()
				delete((*s).clients, client)
			}
		}
	}
}

func (s *Service) ErrorDistribution() {
	for {
		e := <-(*s).errorDistributor
		log.Println(e.Error.Error())
		if (*s).clients[e.Client] {
			s.m.Lock()
			err := e.Client.WriteMessage(websocket.TextMessage, []byte(e.Error.Error()))
			s.m.Unlock()
			if err != nil {
				e.Client.Close()
				delete((*s).clients, e.Client)
			}
		} else {
			log.Println("WTF")
			delete((*s).clients, e.Client)
		}

	}
}
