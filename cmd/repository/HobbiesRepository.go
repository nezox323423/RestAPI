package repository

import (
	"RestAPI/cmd/database"
	"RestAPI/cmd/manager"
	"log"
)

type HobbiesRepository struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	UserId int64  `json:"user_id"`
}

func Hobbies() ([]HobbiesRepository, error) {
	manager.GetEnv()
	db := database.ConnectToMySql()
	defer db.Close()
	query := "SELECT * FROM hobies"
	rows, err := db.Query(query)
	switch true {
	case err != nil:
		log.Fatal(err)
	}
	defer rows.Close()
	var Hobbies []HobbiesRepository
	for rows.Next() {
		var hobby HobbiesRepository
		err := rows.Scan(&hobby.ID, &hobby.Name, &hobby.UserId)
		if err != nil {
			log.Fatal(err)
		}
		Hobbies = append(Hobbies, hobby)
	}

	return Hobbies, nil
}
