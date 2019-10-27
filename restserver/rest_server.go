package restserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GetMethod    = "GET"
	PostMethod   = "POST"
	PutMethod    = "PUT"
	DeleteMethod = "DELETE"
)

// REST server here

type RestServerT interface {
	Start()
	Shutdown()
}

type RestServer struct {
	mux mux.Router
}

var _ RestServerT = (*RestServer)(nil)

func NewRestServer() {
	router := mux.NewRouter()

	// Route bots functions
	botRouter := router.PathPrefix("/bots").Subrouter()
	botRouter.HandleFunc("/", NewBotHandler).Methods(PostMethod)
	botRouter.HandleFunc("/{name}/{channel}", AddBotToChannelHandler).Methods(PostMethod)
	botRouter.HandleFunc("/{name}/{channel}", RemoveBotFromChannelHandler).Methods(DeleteMethod)

	// Route channels functions.
	channelRouter := router.PathPrefix("/channels/{channel}").Subrouter()
	channelRouter.HandleFunc("/", GetChannelHandler).Methods(GetMethod)

	// Route commands functions
	commandRouter := channelRouter.PathPrefix("/commands").Subrouter()
	commandRouter.HandleFunc("/", AddCommandHandler).Methods(PostMethod)
	commandRouter.HandleFunc("/{name}", UpdateCommandHandler).Methods(PutMethod)
	commandRouter.HandleFunc("/{name}", DeleteCommandHandler).Methods(DeleteMethod)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func (server *RestServer) Start() {

}

func (server *RestServer) Shutdown() {

}
