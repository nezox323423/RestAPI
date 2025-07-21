package main

import (
	"RestAPI/cmd/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := router.Router()

	serverErr := make(chan error)

	go func() {
		if err := http.ListenAndServe(":8080", r); err != nil {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("Server error: %v", err)
	case <-stop:
		log.Println("Shutting down gracefully...")
		log.Println("Server stopped")
	}
}
