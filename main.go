package main

import (
	"log"
	"net/http"
	"os"

	"github.com/spf13/viper"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/ojes944/drop1/internal/db"
	"github.com/ojes944/drop1/internal/ws"
)

func main() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("No .env file found or error reading .env: %v", err)
	}

	// Set environment variables from viper
	for _, key := range viper.AllKeys() {
		os.Setenv(key, viper.GetString(key))
	}

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
