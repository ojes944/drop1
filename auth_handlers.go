package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ojes944/drop1/internal/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

var jwtSecret = []byte("supersecret")

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	fmt.Println(req)
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	err := db.CreateUser(req.Email, string(hash), req.Name)
	if err != nil {
		fmt.Println("Error creating user:", err)
		http.Error(w, "user exists", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id, password, name, err := db.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(jwtSecret)
	json.NewEncoder(w).Encode(AuthResponse{Token: tokenStr})
	println(name)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement secure token flow, email sending, and password reset
	w.WriteHeader(http.StatusNotImplemented)
}
