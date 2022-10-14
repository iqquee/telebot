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
			log.Printf("[%s] %s\n", update.Message.From.UserName, update.Message.Text)

			fmt.Printf("This is that message text gottent fron update: %s.........\n", update.Message.Text)

			// THESE APPLIES TO EVERYONE IN THE GROUP EXCEPT THE ADMIN OF THE GROUP
			foundUser := update.Message.From.UserName
			if foundUser != os.Getenv("USER_NAME1") || foundUser != os.Getenv("USER_NAME2") {

				//welcome new users
				wc := update.Message.NewChatMembers
				if wc != nil {
					var byt []map[string]string

					jsonM, _ := json.Marshal(wc)

					json.Unmarshal(jsonM, &byt)
					for val := range byt {
						if byt[val]["first_name"] != "" {
							firstname := byt[val]["first_name"]

							fmt.Printf("This is the firstname found %v.....\n", firstname)
							sendMsg := fmt.Sprintf("Hi %s, Welcome to contact gain with EXODUS kindly add a minimum of 30 telegram user's to be able to drop your number and get verified into the VCF 🗃️ file, if you can't follow the process send a DM to @EXODUSTV96 to pay #500 and get verified to the VCF file", firstname)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendMsg)

							bot.Send(msg)
						}
					}
				}

				if update.Message.NewChatMembers != nil {

					up := update.Message
					var addedUsers database.AddedUsers
					jsonr, _ := json.Marshal(up)

					json.Unmarshal(jsonr, &addedUsers)
					insertID, _ := database.CreateMongoDoc(database.UserCollection, addedUsers)
					fmt.Printf("Mongodb data created with ID: %v\n", insertID)

				}

				usersFirstname := update.Message.From.FirstName
				//delete messages that contains link sent by other users aside from the admin
				domain := urlChecker(update.Message.Text)
				if domain {
					fmt.Println("Message contains a link...")
					deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					bot.Send(deleteMsg)

					fmt.Println("deleted message that contains link...")
					fmt.Println(deleteMsg)
					//notify they user that links can't be sent to the group
					sendMsg := fmt.Sprintf("%s The message you sent contains a link or you failed to add a space after the full stop(.). Links cannot be sent to this group", usersFirstname)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendMsg)

					bot.Send(msg)
				} else { //if the messages sent to the group is not a link

					//check if the text message sent is not empty
					if update.Message.Text != "" {
						//check if the user have already added _ number of users to the group
						countFilter := bson.M{"from.firstname": update.Message.From.FirstName}
						addedUserCount := database.CountCollection(database.UserCollection, countFilter)
						fmt.Printf("This is the number of users you have added to the group %v\n....", addedUserCount)
						userNum := 30
						if addedUserCount < userNum {
							// delete the messages sent to the group by the user who have not added the set numbers of users
							deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
							bot.Send(deleteMsg)

							fmt.Println(deleteMsg)
							// and if not delete their message and notify them to first add _ numbers of users before they can send in messages
							usersToAdd := userNum - addedUserCount
							sendMsg := fmt.Sprintf("%s you have only added %v user(s). You need to add %v more user(s) to be able to send messages to this group. If you can't follow the process send a DM to @EXODUSTV96 to pay #500 and get verified to the VCF file", usersFirstname, addedUserCount, usersToAdd)
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, sendMsg)
							// msg.ReplyToMessageID = update.Message.MessageID
							bot.Send(msg)
						}
					}
				}

			}

		}
	}
}
