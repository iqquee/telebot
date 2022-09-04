package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	// "github.com/gorilla/mux"

	"github.com/gorilla/mux"
	"github.com/hisyntax/telebot/tele"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}
	r := mux.NewRouter()

	for {
		fmt.Printf("for loop running...\n")
		r.HandleFunc("/", RunBot)

		port := os.Getenv("PORT")

		http.ListenAndServe(":"+port, r)
	}
}

func RunBot(w http.ResponseWriter, r *http.Request) {
	tele.Bot()
}
