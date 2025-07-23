package apiserver

import (
	"RestAPI/cmd/exceptions"
	"RestAPI/cmd/repository"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type CreateUserRequestUser struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Age  int32  `json:"age" validate:"required,min=1,max=120"`
}

type CreateUserRequestHobbie struct {
	Name   string `json:"name" validate:"required,min=2,max=50"`
	UserID int64  `json:"user_id" validate:"min=1,max=100"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		exceptions.ValidateIdRequest(w)
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
		exceptions.NotExistInDb(w)
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

	var req CreateUserRequestUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат запроса",
		})
		return
	}
	// Валидация данных
	if req.Name == "" {
		exceptions.ValidateNameRequest(w)
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
		exceptions.ValidateIdRequest(w)
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
		exceptions.NotExistInDb(w)
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

func GetHobbies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	hobbies, _ := repository.Hobbies()
	if len(hobbies) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Нет ни одного хобби",
		})
	}
	if len(hobbies) != 0 {
		if err := json.NewEncoder(w).Encode(hobbies); err != nil && len(hobbies) != 0 {
			log.Printf("JSON encoding error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func GetHobbie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		exceptions.ValidateIdRequest(w)
		return
	}
	hobbie, exist, err := repository.GetHobbieById(id)
	if !exist {
		exceptions.NotExistInDb(w)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	if err := json.NewEncoder(w).Encode(hobbie); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func CreateHobbies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CreateUserRequestHobbie
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Данные не корректные",
		})
		return
	}
	if req.Name == "" {
		exceptions.ValidateNameRequest(w)
		return
	}

	var hobbie repository.HobbiesRepository
	var err error
	var exist bool

	if req.UserID == 0 {
		hobbie, err, exist = repository.CreateHobbie(req.Name, nil)
	} else {
		hobbie, err, exist = repository.CreateHobbie(req.Name, &req.UserID)
	}
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "С такким user_id нет хобби",
		})
		return
	}
	if err != nil {
		exceptions.NotImplemented(w)
		return
	}
	json.NewEncoder(w).Encode(hobbie)
}
