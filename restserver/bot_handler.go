package restserver

import "net/http"

//botRouter.HandleFunc("/", NewBotHandler).Methods(PostMethod)
func NewBotHandler(w http.ResponseWriter, r *http.Request) {
	return
}

//botRouter.HandleFunc("/{name}/{channel}", AddBotToChannelHandler).Methods(PostMethod)
func AddBotToChannelHandler(w http.ResponseWriter, r *http.Request) {
	return
}

//botRouter.HandleFunc("/{name}/{channel}", RemoveBotFromChannelHandler).Methods(DeleteMethod)
func RemoveBotFromChannelHandler(w http.ResponseWriter, r *http.Request) {
	return
}
