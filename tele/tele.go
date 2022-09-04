package tele

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hisyntax/telebot/database"
	"go.mongodb.org/mongo-driver/bson"
)

func urlChecker(character string) bool {
	var val bool
	//check if the string contains a .
	if strings.Contains(character, ".") {
		for i, value := range character {
			//get the index of .
			if string(value) == "." {
				fmt.Printf("%d - %v\n", i, string(value))
				prev := i - 1
				next := i + 1
				for e, v := range character {
					//check the previous character if its an "" string
					if e == prev {
						if string(v) != " " && string(v) != "." {
							//check the next character if its an "" string
							for ee, vv := range character {
								if ee == next {
									if string(vv) != " " && string(vv) != "." {
										val = true
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return val
}

func Bot() {
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// THESE APPLIES TO EVERYONE IN THE GROUP EXCEPT THE ADMIN OF THE GROUP
			foundUser := update.Message.From.UserName
			if foundUser != os.Getenv("USER_NAME") {

				//welcome new users
				wc := update.Message.NewChatMembers
				if wc != nil {
					var byt []map[string]string

					jsonM, _ := json.Marshal(wc)

					json.Unmarshal(jsonM, &byt)
					for val := range byt {
						if byt[val]["username"] != "" {
							username := byt[val]["username"]
							sendMsg := fmt.Sprintf("@%s welcome to test-bot group, we catch fun here :)", username)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendMsg)

							bot.Send(msg)
						}
					}

					// check if the user was already added to the group before

				}

				// if update.Message != nil {
				//check if user is not removed from the group and user is not added to the group
				//every other messages should be deleted

				// removeUser := update.Message.LeftChatMember.UserName // get the removed user
				// var addUser bool
				// fmt.Printf("%v.....\n", removeUser)
				// fmt.Printf("%v.....\n", addUser)

				// addUserUpdate := update.Message.NewChatMembers
				// for i, v := range addUserUpdate {
				// 	if i == 0 {
				// 		fmt.Println(v)
				// 		addUser = true
				// 	}
				// }
				// var addedUser []map[string]string
				// aUserJson, _ := json.Marshal(addUserUpdate)
				// json.Unmarshal(aUserJson, &addedUser)
				// for val := range addedUser {
				// 	username := addedUser[val]["username"]
				// 	addUser = username
				// }

				// 	if removeUser == "" {
				// 		// check the database to get the number of users a particular user have added
				// 		// to the group to know if they are eligible to send messages to the group in this case 30
				// 		countFilter := bson.M{"from.username": update.Message.From.UserName}
				// 		addedUserCount := database.CountCollection(database.UserCollection, countFilter)
				// 		if addedUserCount < 2 {
				// 			//delete the messages sent to the group by the user
				// 			deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
				// 			bot.Send(deleteMsg)

				// 			fmt.Println(deleteMsg)
				// 			//notify the users that that they need to add 30 people to the group
				// 			sendMsg := fmt.Sprintf("@%s you need to add 30 users to be able to send messages to this group", foundUser)
				// 			msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendMsg)
				// 			// msg.ReplyToMessageID = update.Message.MessageID
				// 			bot.Send(msg)
				// 		}
				// 	}
				// }

				//check if the recently added user already exists in the datatase
				filter := bson.M{"from.username": update.Message.From.UserName}
				res, _ := database.GetMongoDoc(database.UserCollection, filter)
				if res != nil {
					//loop through the object to get the username of the just added user
					var byt []map[string]interface{}
					arr := res.NewChatMembers
					jm, _ := json.Marshal(arr)

					json.Unmarshal(jm, &byt)

					//add _ number users to the group before being able to send messages to the group
					up := update.Message
					var addedUsers database.AddedUsers
					jsonr, _ := json.Marshal(up)

					json.Unmarshal(jsonr, &addedUsers)

					var newUser []map[string]interface{}
					addedUserR := addedUsers.NewChatMembers
					addedUserjson, _ := json.Marshal(addedUserR)

					json.Unmarshal(addedUserjson, &newUser)

					//so long newchatmembers object is not nil - it would ignore updated when a user is removed from the group
					if addedUsers.NewChatMembers != nil {
						for v := range newUser {
							if newUser[v]["username"] != "" {
								for vv := range byt {
									//check if the user have not been added to the group before
									if byt[vv]["username"] != newUser[v]["username"] {
										//add the new user
										insertID, _ := database.CreateMongoDoc(database.UserCollection, addedUsers)
										// if err != nil {
										// 	fmt.Printf("Mongo db error: %v\n", err)
										// }
										fmt.Printf("Mongodb data created with ID: %v\n", insertID)
									}
								}
							}
						}
					}
				} else {
					up := update.Message
					var addedUsers database.AddedUsers
					jsonr, _ := json.Marshal(up)

					json.Unmarshal(jsonr, &addedUsers)
					insertID, _ := database.CreateMongoDoc(database.UserCollection, addedUsers)
					fmt.Printf("Mongodb data created with ID: %v\n", insertID)
				}

				//delete messages that contains link sent by other users aside from the admin
				domain := urlChecker(update.Message.Text)
				if domain {
					fmt.Println("Message contains a link...")
					deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					bot.Send(deleteMsg)
					fmt.Println("deleted message that contains link...")
					fmt.Println(deleteMsg)
				}
			}

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			// bot.Send(msg)
		}
	}

	mux := http.NewServeMux()
	port := os.Getenv("PORT")
	go http.ListenAndServe(":"+port, mux)
}
