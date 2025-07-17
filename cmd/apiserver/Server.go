package apiserver

import (
	"RestAPI/cmd/repository"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Age  int32  `json:"age" validate:"required,min=1,max=120"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат ID",
		})
		return
	}
	user, err := repository.DeleteUser(id)
	switch true {
	case user == "deleted user":
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Вы удалили пользователя",
		})
	case err != nil:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "<UNK> <UNK>",
		})
	case user:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователя с таким id не существует",
		})
	default:
		json.NewEncoder(w).Encode(user)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "<UNK> <UNK>",
		})
	}

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат запроса",
		})
		return
	}
	// Валидация данных
	if req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Имя не может быть пустым",
		})
		return
	}

	if req.Age <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Возраст должен быть положительным числом",
		})
		return
	}
	user, err := repository.CreateUserInDb(req.Name, req.Age)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "<UNK> <UNK> <UNK>",
		})
		return
	}
	if user == true {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Такой пользователь уже есть",
		})
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат ID",
		})
		return
	}
	user, exists := repository.GetUserById(id)

	switch exists {
	case true:
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("Ошибка кодирования пользователя: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	case false:
		// Если пользователь не найден
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{
			"error": fmt.Sprintf("пользователь с ID %d не найден", id),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Ошибка кодирования ошибки: %v", err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}

	}

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users := repository.Users()
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
