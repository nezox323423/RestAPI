package repository

import (
	"RestAPI/cmd/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

// User представляет модель пользователя
type User struct {
	ID   int32  `json:"id"`
	Age  int32  `json:"age"`
	Name string `json:"name"`
}

type UserStore interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id int32) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id int32) error
}

type MySQLUserStore struct {
	db database.DBConnection
}

// NewMySQLUserStore создает новый экземпляр MySQLUserStore
func NewMySQLUserStore(conn database.DBConnection) *MySQLUserStore {
	return &MySQLUserStore{db: conn}
}

func (s *MySQLUserStore) GetAll(ctx context.Context) ([]User, error) {
	query := "SELECT id, age, name FROM users"
	rows, err := s.db.GetDB().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Age, &u.Name); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

func (s *MySQLUserStore) GetByID(ctx context.Context, id int32) (*User, error) {
	query := "SELECT id, age, name FROM users WHERE id = ?"
	row := s.db.GetDB().QueryRowContext(ctx, query, id)

	var u User
	err := row.Scan(&u.ID, &u.Age, &u.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &u, nil
}

func (s *MySQLUserStore) Create(ctx context.Context, user *User) (*User, error) {
	// Проверяем, существует ли уже такой пользователь
	existingUser, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return existingUser, nil
	}

	query := "INSERT INTO users (name, age) VALUES (?, ?)"
	result, err := s.db.GetDB().ExecContext(ctx, query, user.Name, user.Age)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return s.GetByID(ctx, int32(id))
}

func (s *MySQLUserStore) Delete(ctx context.Context, id int32) error {
	// Проверяем существование пользователя
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	query := "DELETE FROM users WHERE id = ?"
	_, err = s.db.GetDB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
