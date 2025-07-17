package main

import (
	"RestAPI/cmd/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	log.Fatal(http.ListenAndServe(":8080", r))
}
