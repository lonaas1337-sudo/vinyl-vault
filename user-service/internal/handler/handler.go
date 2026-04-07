package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/model"
	"github.com/lonaas1337-sudo/vinylvault/user-service/internal/repository"
)

var repo *repository.UserRepository

func SetRepository(r *repository.UserRepository) {
	repo = r
}

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Status:  "healthy",
		Service: "user",
		Version: "0.1.0",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Received user registration for email %s", req.Email)

	user, err := model.NewUser(req.Email, req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, err := repo.CreateUser(r.Context(), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := RegisterResponse{
		Message: "user registered successfully",
		ID:      userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
