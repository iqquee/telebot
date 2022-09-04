package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DbInstance() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	MongoDb := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully connected to mongodb")
	return client
}

var Client *mongo.Client = DbInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
	databaseName := os.Getenv("DATABASE_NAME")
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)
	return collection
}

func GetMongoDoc(colName *mongo.Collection, filter interface{}) (*AddedUsers, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data AddedUsers

	if err := colName.FindOne(ctx, filter).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func CreateMongoDoc(colName *mongo.Collection, data interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertNum, insertErr := colName.InsertOne(ctx, data)
	if insertErr != nil {
		return nil, insertErr
	}

	return insertNum, nil
}

func CountCollection(colName *mongo.Collection, filter interface{}) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := colName.CountDocuments(ctx, filter)

	if err != nil {
		return 0
	}
	return int(count)
}
