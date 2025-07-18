package exceptions

import (
	"encoding/json"
	"net/http"
)

func ValidateIdRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Неверный формат ID",
	})
}
func ValidateNameRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Имя не должно быть пустым",
	})
}

func NotExistInDb(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "С таким id не существует",
	})
}
