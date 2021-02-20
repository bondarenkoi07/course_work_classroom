package router

import (
	"awesomeProject1/Service"
	goji "goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

type Router struct {
	mux     *goji.Mux
	service Service.Service
}

func NewRouter() *Router {
	router := &Router{
		mux:     goji.NewMux(),
		service: Service.NewService(),
	}
	go router.service.ClassroomControllerDistribution()
	go router.service.ErrorDistribution()
	return router
}

func (r *Router) AddIndex() {
	(*r).mux.Handle(pat.Get("/*"), (*r).service.Index())
}

func (r *Router) AddAction() {
	(*r).mux.HandleFunc(pat.Get("/action"), (*r).service.Action)
}

func (r *Router) Exec() {
	err := http.ListenAndServe(":8080", (*r).mux)
	if err != nil {
		log.Fatal(err)
	}
}
