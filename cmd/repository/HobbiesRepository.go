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
	UserId *int64 `json:"user_id"`
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

func CreateHobbie(name string, userId *int64) (HobbiesRepository, error, bool) {
	manager.GetEnv()
	db := database.ConnectToMySql()
	defer db.Close()

	var query string
	var result sql.Result
	var err error
	var hobby HobbiesRepository

	if userId == nil {
		query = "INSERT INTO hobies (name) VALUES (?)"
		result, err = db.Exec(query, name)
	} else {
		query = "select name from users where id = ?"
		row := db.QueryRow(query, userId)
		var nameUser string
		if err := row.Scan(&nameUser); err != nil {
			if err == sql.ErrNoRows {
				return hobby, err, false
			}
		}
		query = "INSERT INTO hobies (name,user_id) VALUES (?,?)"
		result, err = db.Exec(query, name, userId)
	}
	if err != nil {
		log.Fatal(err)
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	queryGet := "SELECT id, name, user_id FROM hobies WHERE id = ?"
	row := db.QueryRow(queryGet, insertId)

	err = row.Scan(&hobby.ID, &hobby.Name, &hobby.UserId)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}
	return hobby, nil, true
}
