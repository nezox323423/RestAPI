package database

import (
	"database/sql"
	"log"
	"os"
)

type MySQLConfig struct {
	User string
	Pass string
	Name string
}

func NewMySQLConfig() *MySQLConfig {
	return &MySQLConfig{
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Name: os.Getenv("DB_NAME"),
	}
}

func (c *MySQLConfig) BuildDSN() string {
	return c.User + ":" + c.Pass + "@/" + c.Name
}

type MySQLConnection struct {
	db *sql.DB
}

func (conn *MySQLConnection) Connect(dsn string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	conn.db = db
	return nil
}

func (conn *MySQLConnection) GetDB() *sql.DB {
	return conn.db
}

func NewMySQLConnection() (*MySQLConnection, error) {
	config := NewMySQLConfig()
	conn := &MySQLConnection{}
	err := conn.Connect(config.BuildDSN())
	if err != nil {
		return nil, err
	}
	log.Println("Connected to DB")
	return conn, nil
}
