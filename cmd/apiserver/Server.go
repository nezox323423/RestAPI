package apiserver

import (
	"RestAPI/cmd/repository"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type APIServer struct {
	userStore  repository.UserStore
	hobbyStore repository.HobbiesStore
	router     *mux.Router
}

func NewAPIServer(userStore repository.UserStore, hobbyStore repository.HobbiesStore) *APIServer {
	s := &APIServer{
		userStore:  userStore,
		hobbyStore: hobbyStore,
		router:     mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/users", s.handleGetUsers).Methods("GET")
	s.router.HandleFunc("/users/{id}", s.handleGetUser).Methods("GET")
	s.router.HandleFunc("/users", s.handleCreateUser).Methods("POST")
	s.router.HandleFunc("/users/{id}", s.handleDeleteUser).Methods("DELETE")

	s.router.HandleFunc("/hobbies", s.handleGetHobbies).Methods("GET")
	s.router.HandleFunc("/hobbies/{id}", s.handleGetHobby).Methods("GET")
	s.router.HandleFunc("/hobbies", s.handleCreateHobby).Methods("POST")
}

// DTO (Data Transfer Objects)
type (
	CreateUserRequest struct {
		Name string `json:"name" validate:"required,min=2,max=50"`
		Age  int32  `json:"age" validate:"required,min=1,max=120"`
	}

	CreateHobbyRequest struct {
		Name   string `json:"name" validate:"required,min=2,max=50"`
		UserID *int64 `json:"user_id" validate:"omitempty,min=1"`
	}
)

// Обработчики запросов
func (s *APIServer) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := s.userStore.GetAll(ctx)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get users")
		return
	}
	s.respondWithJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := s.userStore.GetByID(ctx, int32(id))
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get user")
		return
	}
	if user == nil {
		s.respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	s.respondWithJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user := &repository.User{
		Name: req.Name,
		Age:  req.Age,
	}

	createdUser, err := s.userStore.Create(ctx, user)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	s.respondWithJSON(w, http.StatusCreated, createdUser)
}

func (s *APIServer) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := s.userStore.Delete(ctx, int32(id)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.respondWithError(w, http.StatusNotFound, "user not found")
			return
		}
		s.respondWithError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	s.respondWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}

func (s *APIServer) handleGetHobbies(w http.ResponseWriter, r *http.Request) {
	hobbies, err := s.hobbyStore.GetAll()
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get hobbies")
		return
	}
	if len(hobbies) == 0 {
		s.respondWithJSON(w, http.StatusOK, []interface{}{})
		return
	}
	s.respondWithJSON(w, http.StatusOK, hobbies)
}

func (s *APIServer) handleGetHobby(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid hobby ID")
		return
	}

	hobby, err := s.hobbyStore.GetByID(id)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to get hobby")
		return
	}
	if hobby == nil {
		s.respondWithError(w, http.StatusNotFound, "hobby not found")
		return
	}

	s.respondWithJSON(w, http.StatusOK, hobby)
}

func (s *APIServer) handleCreateHobby(w http.ResponseWriter, r *http.Request) {
	var req CreateHobbyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	hobby := &repository.HobbiesRepository{
		Name:   req.Name,
		UserId: req.UserID,
	}

	createdHobby, err := s.hobbyStore.Create(hobby)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.respondWithError(w, http.StatusBadRequest, "user not found")
			return
		}
		s.respondWithError(w, http.StatusInternalServerError, "failed to create hobby")
		return
	}

	s.respondWithJSON(w, http.StatusCreated, createdHobby)
}

// Вспомогательные методы
func (s *APIServer) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, map[string]string{"error": message})
}

func (s *APIServer) respondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
