package tele

import (
	"encoding/json"
	"fmt"
	"log"
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
			up := update.Message
			fmt.Printf("This is the received message... %v\n", up)

			//welcome new users
			wc := update.Message.NewChatMembers
			if wc != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "welcome to test-bot group. Please 30 people to be able to send messages in the group")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)

				//new users should not be able to send messages to the group until they add 30 more persons to the group

				// check if the user was already added to the group before

			}

			//check if the recently added user already exists in the datatase
			filter := bson.M{"from.username": update.Message.From.UserName}
			res, _ := database.GetMongoDoc(database.UserCollection, filter)

			//loop through the object to get the username of the just added user
			var byt []map[string]interface{}
			arr := res.NewChatMembers
			jm, _ := json.Marshal(arr)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			json.Unmarshal(jm, &byt)
			// ; err != nil {
			// 	fmt.Println(err)
			// }

			//add _ number users to the group before being able to send messages to the group
			var addedUsers database.AddedUsers
			jsonr, _ := json.Marshal(up)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			json.Unmarshal(jsonr, &addedUsers)
			// ; err != nil {
			// 	fmt.Println(err)
			// }

			var newUser []map[string]interface{}
			addedUserR := addedUsers.NewChatMembers
			addedUserjson, _ := json.Marshal(addedUserR)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			json.Unmarshal(addedUserjson, &newUser)
			//  err != nil {
			// 	fmt.Println(err)
			// }

			//so long newchatmembers object is not nil - it would ignore updated when a user is removed from the group
			if addedUsers.NewChatMembers != nil {
				for _, v := range newUser {
					if v["username"] != "" {
						for _, vv := range byt {
							//check if the user have not been added to the group before
							if vv["username"] != v["username"] {
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

			//delete messages that contains link sent by other users aside from the admin
			adminUser := update.Message.From.UserName
			if adminUser != os.Getenv("USER_NAME") {
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
}
