package router

import (
	"RestAPI/cmd/apiserver"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", apiserver.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", apiserver.GetUser).Methods("GET")
	r.HandleFunc("/users", apiserver.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", apiserver.DeleteUser).Methods("DELETE")
	return r
}
