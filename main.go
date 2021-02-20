package main

import (
	"awesomeProject1/router"
	_ "github.com/gorilla/websocket"
)

func main() {
	Router := router.NewRouter()
	Router.AddAction()
	//Router.AddIndex() - old version, now static is served by Nginx
	Router.Exec()
}
