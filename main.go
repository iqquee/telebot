package main

import (
	"log"

	// "github.com/gorilla/mux"

	"github.com/hisyntax/telebot/tele"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	tele.Bot()
	// r := mux.NewRouter()

	// fmt.Printf("for loop running...\n")
	// r.HandleFunc("/", RunBot)

	// port := os.Getenv("PORT")

	// http.ListenAndServe(":"+port, r)

}

// func RunBot(w http.ResponseWriter, r *http.Request) {

// }
