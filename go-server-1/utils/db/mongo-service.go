package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Config struct {
	Uri      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (mongoConfig *Config) Connect() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoConfig.Uri)
	clientOptions.SetAuth(options.Credential{
		Username: mongoConfig.Username,
		Password: mongoConfig.Password,
	})
	mongodbConnection, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error Connecting to MongoDB: ", err)
		return nil, err
	}
	fmt.Println("Mongo Connection Created")
	return mongodbConnection, nil
}
