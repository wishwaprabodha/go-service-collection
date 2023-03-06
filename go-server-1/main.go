package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/router"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/service"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/service/handler"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/utils/amqp/rmq"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/utils/db"
	"log"
	"net/http"
)

func main() {
	log.Println("Connecting To RMQ...")

	rmqConfigs := rmq.NewQueueHandler("rmq-q", "rmq-rk", "rmq-ex", 5, 3)
	// create queue 01
	//userQueueConfig := rmq.NewQueueHandler("user-q", "user-rk", "user-ex", 5, 3)
	rmqConn, rmqErr := rmqConfigs.Connection("amqp://guest:guest@rabbitmq/")
	//userRmqConn, userRmqErr = userQueueConfig.Connection("amqp://guest:guest@rabbitmq/")
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
	_ = db.DbConnection()
	log.Println("MySQL Connection Successful")
	_, detaErr := db.DetaConnection()
	if detaErr != nil {
		log.Panic("Deta Connection Error: ", detaErr)
	}
	log.Println("Deta Connection Success")
	service.InitBooks()
	rmqChannel, rmqChErr := rmqConfigs.CreateQueue(rmqConn)
	if rmqChErr != nil {
		log.Panic("RMQ Channel Creation Error: ", rmqChErr)
	}
	log.Println("Connected To RMQ Channel ...")
	rmqPubErr := rmqConfigs.PublishMessage(rmqChannel, "Hello World")
	if rmqPubErr != nil {
		log.Panic("RMQ Channel Publish Error: ", rmqPubErr)
	}
	log.Println("Message Published")
	userEventHandler := handler.UserHandler{
		Config: rmqConfigs,
	}
	rmqConsumeErr := rmqConfigs.ConsumeMessage(context.Background(), rmqChannel, userEventHandler)
	if rmqConsumeErr != nil {
		log.Panic("RMQ Channel Consume Error: ", rmqConsumeErr)
	}
	log.Println(http.ListenAndServe(":8001", router.StartRoutes()))
	log.Println("Server Started")
}
