package database

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	UserCollection  *mongo.Collection = OpenCollection(Client, "users")
	ErrUsersLess                      = errors.New("the users you added is less than 30")
	ErrUsersLessMSg                   = errors.New("please add up to 30 users before you can semd messages to this group")
)

type AddedUsers struct {
	From           From             `json:"from"`
	NewChatMembers []NewChatMembers `json:"new_chat_members"`
}

type From struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type NewChatMembers struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}
