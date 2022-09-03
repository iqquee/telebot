package main

import (
	"log"

	// "github.com/hisyntax/telebot/tele"
	// "github.com/gorilla/mux"

	"github.com/hisyntax/telebot/tele"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	tele.Bot()
}
