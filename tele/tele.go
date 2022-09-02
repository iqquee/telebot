package tele

import (
	"fmt"
	"log"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			wc := update.Message.NewChatMembers
			if wc != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "welcome to test-bot group")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}

			//check the role of user who sent link if admin else delete the message
			adminUser := update.Message.From.UserName
			if adminUser != os.Getenv("USER_NAME") {
				_, err := url.ParseRequestURI(update.Message.Text)
				if err != nil {
					fmt.Println("Message contains a link...")
					deleteMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
					bot.Send(deleteMsg)
					fmt.Println("deleting message that contains link...")
					fmt.Println(deleteMsg)
				}
			}

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			// bot.Send(msg)
		}
	}
}
