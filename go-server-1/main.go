package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/wishwaprabodha/go-server/router"
	"github.com/wishwaprabodha/go-server/service/models"
	"github.com/wishwaprabodha/go-server/utils/amqp"
	"github.com/wishwaprabodha/go-server/utils/db"
	"log"
	"net/http"
)

func main() {
	log.Println("Connecting To RMQ...")
	rmqConn, rmqErr := amqp.Connection()
	if rmqErr != nil {
		log.Panic("RMQ Connection Error: ", rmqErr)
	}
	log.Println("Connected To RMQ...")
	mongoConfig := db.Config{
		Uri:      "mongodb://127.0.0.1:27017",
		Username: "root",
		Password: "password",
	}
	_, err := mongoConfig.Connect()
	if err != nil {
		log.Panic("Mongo Connection Error: ", err)
	}
	_ = db.DBConnection()
	log.Println("MySQL Connection Successful")
	_, detaErr := db.DetaConnection()
	if detaErr != nil {
		log.Panic("Deta Connection Error: ", detaErr)
	}
	log.Println("Deta Connection Success")
	models.InitBooks()
	rmqChannel, rmqChErr := amqp.CreateQueue(rmqConn)
	if rmqChErr != nil {
		log.Panic("RMQ Channel Creation Error: ", rmqChErr)
	}
	log.Println("Connected To RMQ Channel ...")
	rmqPubErr := amqp.PublishMessage(rmqChannel)
	if rmqPubErr != nil {
		log.Panic("RMQ Channel Publish Error: ", rmqPubErr)
	}
	log.Println("Message Published")
	rmqConsumeErr := amqp.ConsumeMessage(rmqChannel)
	if rmqConsumeErr != nil {
		log.Panic("RMQ Channel Consume Error: ", rmqConsumeErr)
	}
	log.Println(http.ListenAndServe(":8001", router.StartRoutes()))
	log.Println("Server Started")
}
