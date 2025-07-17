package database

import (
	"database/sql"
	"log"
	"os"
)

func ConnectToMySql() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to DB")
	}
	return db
}
