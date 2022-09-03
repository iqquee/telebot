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
