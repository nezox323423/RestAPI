package repository

import (
	"RestAPI/cmd/database"
	"database/sql"
	"errors"
	"log"
)

type HobbiesRepository struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	UserId *int64 `json:"user_id"`
}

type HobbiesStore interface {
	GetAll() ([]HobbiesRepository, error)
	GetByID(id int64) (*HobbiesRepository, error)
	Create(hobby *HobbiesRepository) (*HobbiesRepository, error)
}

type MySQLHobbiesStore struct {
	db database.DBConnection
}

func NewMySQLHobbiesStore(conn database.DBConnection) *MySQLHobbiesStore {
	return &MySQLHobbiesStore{db: conn}
}

func (s *MySQLHobbiesStore) GetAll() ([]HobbiesRepository, error) {
	query := "SELECT id, name, user_id FROM hobbies"
	rows, err := s.db.GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hobbies []HobbiesRepository
	for rows.Next() {
		var hobby HobbiesRepository
		if err := rows.Scan(&hobby.ID, &hobby.Name, &hobby.UserId); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		hobbies = append(hobbies, hobby)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return hobbies, nil
}

func (s *MySQLHobbiesStore) GetByID(id int64) (*HobbiesRepository, error) {
	query := "SELECT id, name, user_id FROM hobbies WHERE id = ?"
	row := s.db.GetDB().QueryRow(query, id)

	var hobby HobbiesRepository
	err := row.Scan(&hobby.ID, &hobby.Name, &hobby.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &hobby, nil
}

func (s *MySQLHobbiesStore) Create(hobby *HobbiesRepository) (*HobbiesRepository, error) {
	var query string
	var args []interface{}

	if hobby.UserId == nil {
		query = "INSERT INTO hobbies (name) VALUES (?)"
		args = []interface{}{hobby.Name}
	} else {
		// Проверяем существование пользователя
		userExists, err := s.userExists(*hobby.UserId)
		if err != nil {
			return nil, err
		}
		if !userExists {
			return nil, errors.New("user does not exist")
		}

		query = "INSERT INTO hobbies (name, user_id) VALUES (?, ?)"
		args = []interface{}{hobby.Name, *hobby.UserId}
	}

	result, err := s.db.GetDB().Exec(query, args...)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetByID(id)
}

func (s *MySQLHobbiesStore) userExists(userID int64) (bool, error) {
	query := "SELECT 1 FROM users WHERE id = ? LIMIT 1"
	row := s.db.GetDB().QueryRow(query, userID)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
