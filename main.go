package main

import (
	"RestAPI/cmd/apiserver"
	"RestAPI/cmd/database"
	"RestAPI/cmd/manager"
	"RestAPI/cmd/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	manager.GetEnv()
	conn, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatal(err)
	}
	db := conn.GetDB()
	defer db.Close()

	userStore := repository.NewMySQLUserStore(conn)
	hobbyStore := repository.NewMySQLHobbiesStore(conn)

	api := apiserver.NewAPIServer(userStore, hobbyStore)

	serverErr := make(chan error)

	go func() {
		if err := http.ListenAndServe(":8080", api); err != nil {
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
