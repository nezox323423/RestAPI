package router

import (
	"RestAPI/cmd/apiServer"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", apiServer.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", apiServer.GetUser).Methods("GET")
	r.HandleFunc("/users", apiServer.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", apiServer.DeleteUser).Methods("DELETE")
	return r
}
