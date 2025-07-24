package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DBConnection interface {
	GetDB() *sql.DB
	Ping() error
	Close() error
}

func (conn *MySQLConnection) Ping() error {
	return conn.db.Ping()
}

func (conn *MySQLConnection) Close() error {
	return conn.db.Close()
}
