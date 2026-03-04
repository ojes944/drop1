package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"drop1/internal/db"
	"drop1/internal/ws"
)

func main() {
	db.Init()
	db.InitRedis()

	r := mux.NewRouter()
	r.HandleFunc("/api/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/reset_password", ResetPasswordHandler).Methods("POST")
	r.HandleFunc("/ws", ws.WebSocketHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("web-client/public")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
