package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type user struct {
	Id   int32
	Age  int32
	Name string
}

func Users() []user {
	getEnv()
	db := ConnectToMySql()
	defer db.Close()
	query := "select * from users"
	rows, err := db.Query(query)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	users := []user{}
	for rows.Next() {
		u := user{}
		if err := rows.Scan(&u.Id, &u.Age, &u.Name); err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	return users
}

func GetUserById(id int64) (user, bool) {
	getEnv()
	db := ConnectToMySql()
	defer db.Close()

	query := "SELECT id, age, name FROM users WHERE id = ?"
	row := db.QueryRow(query, id)

	var u user
	err := row.Scan(&u.Id, &u.Age, &u.Name)

	if err != nil {
		if err == sql.ErrNoRows {

			return user{}, false
		}
		log.Printf("Database error: %v", err)
		return user{}, false
	}

	return u, true
}

func CreateUserInDb(name string, age int32) (interface{}, error) {
	getEnv()
	db := ConnectToMySql()
	defer db.Close()

	querySel := "SELECT id, name, age FROM users WHERE name = ? AND age = ?"
	row := db.QueryRow(querySel, name, age)

	var u user
	err := row.Scan(&u.Id, &u.Name, &u.Age)

	if err == nil {
		return true, nil
	}

	if err != sql.ErrNoRows {
		log.Printf("Database select error: %v", err)
		return user{}, fmt.Errorf("database error: %w", err)
	}

	queryIns := "INSERT INTO users (name, age) VALUES (?, ?)"
	result, err := db.Exec(queryIns, name, age)
	if err != nil {
		log.Printf("Database insert error: %v", err)
		return user{}, fmt.Errorf("failed to insert user: %w", err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		log.Printf("Database last insert ID error: %v", err)
		return user{}, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	queryGet := "SELECT id, name, age FROM users WHERE id = ?"
	row = db.QueryRow(queryGet, insertId)

	var newUser user
	if err := row.Scan(&newUser.Id, &newUser.Name, &newUser.Age); err != nil {
		log.Printf("Database select new user error: %v", err)
		return user{}, fmt.Errorf("failed to fetch created user: %w", err)
	}

	return newUser, nil
}

func getEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

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

func DeleteUser(id int64) (interface{}, error) {
	getEnv()
	db := ConnectToMySql()
	defer db.Close()

	println(id)

	querySel := "Select id, name, age FROM users WHERE id = ?"
	row := db.QueryRow(querySel, id)
	var u user
	err := row.Scan(&u.Id, &u.Name, &u.Age)
	switch true {
	case err == sql.ErrNoRows:
		log.Printf("Пользователя нет")
		return true, nil
	case err != nil:
		log.Printf("Database last insert ID error: %v", err)
		return user{}, fmt.Errorf("failed to get last insert ID: %w", err)
		//case err == nil:
		//	return true, nil
	}
	query := "DELETE FROM users WHERE id = ?"
	_, err = db.Exec(query, id)
	if err != nil {
		log.Printf("Database delete user error: %v", err)
	}
	return "deleted user", nil
}
