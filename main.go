package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/hisyntax/telebot/tele"
	// "github.com/gorilla/mux"
	"github.com/hisyntax/telebot/database"
	"github.com/hisyntax/telebot/tele"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	// r := mux.NewRouter()

	// r.HandleFunc("/docs", GetDocs).Methods("GET")

	// fmt.Println("server is starting...")
	// if err := http.ListenAndServe(":8080", r); err != nil {
	// 	fmt.Println(err)
	// }
	tele.Bot()

}

func GetDocs(w http.ResponseWriter, r *http.Request) {
	filter := bson.M{"from.lastname": "Jiya"}
	res, err := database.GetMongoDoc(database.UserCollection, filter)
	if err != nil {
		fmt.Println(err)
	}

	// for i, v := range res.NewChatMembers {

	var byt []map[string]interface{}
	arr := res.NewChatMembers
	jsonr, err := json.Marshal(arr)
	if err != nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(jsonr, &byt); err != nil {
		fmt.Println(err)
	}

	for _, v := range byt {
		if v["username"] == "samurai1979" {
			fmt.Println(v["username"])
			json.NewEncoder(w).Encode(v["username"])
		}
	}
	// fmt.Printf("%d =  %v\n", i, v)

	// }

	fmt.Printf("This is the response from database: %v\n", res)

	// json.NewEncoder(w).Encode(mapping)
}
