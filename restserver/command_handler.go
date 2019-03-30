package restserver

import "net/http"

//commandRouter.HandleFunc("/", AddCommandHandler).Methods(PostMethod)
func AddCommandHandler(w http.ResponseWriter, r *http.Request) {
	return
}

//commandRouter.HandleFunc("/{name}", UpdateCommandHandler).Methods(PutMethod)
func UpdateCommandHandler(w http.ResponseWriter, r *http.Request) {
	return
}

// commandRouter.HandleFunc("/{name}", DeleteCommandHandler).Methods(DeleteMethod)
func DeleteCommandHandler(w http.ResponseWriter, r *http.Request) {
	return
}
