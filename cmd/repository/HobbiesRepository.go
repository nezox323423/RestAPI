package repository

import (
	"RestAPI/cmd/database"
	"RestAPI/cmd/manager"
	"database/sql"
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

func GetHobbieById(id int64) (HobbiesRepository, bool, error) {
	manager.GetEnv()
	db := database.ConnectToMySql()
	defer db.Close()
	query := "SELECT * FROM hobies WHERE id=?"
	row := db.QueryRow(query, id)

	var hobby HobbiesRepository

	err := row.Scan(&hobby.ID, &hobby.Name, &hobby.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return hobby, false, err
		}
		log.Fatal("Database error: %v", err)
	}
	return hobby, true, nil
}

func CreateHobbie(name string, userId *int64) (HobbiesRepository, error) {
	//todo нужно подумать, как валидировать если запись уже в бд есть.
	//todo У нас в целом нет никаих уникальных данных, которые не могут
	//todo повторяться в других записях
	manager.GetEnv()
	db := database.ConnectToMySql()
	defer db.Close()
	switch userId {
	case nil:
		query := "INSERT INTO hobies (name) VALUES (?)"
	default:

	}
}
