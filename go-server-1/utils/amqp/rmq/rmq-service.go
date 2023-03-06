package rmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/service/handler"
	"log"
	"math"
	"time"
)

type Config struct {
	qName      string
	rKey       string
	exchange   string
	retryCount int64
	interval   time.Duration
}

type QueueService interface {
	Connection(url string) (*amqp.Connection, error)
	CreateQueue(connection *amqp.Connection) (*amqp.Channel, error)
	PublishMessage(channel *amqp.Channel) error
	ConsumeMessage(channel *amqp.Channel) error
	Reconnect(conn *amqp.Connection)
}

func NewQueueHandler(qName string, routineKey string, exchange string, retryCount int64, interval time.Duration) *Config {
	return &Config{qName: qName, rKey: routineKey, exchange: exchange, retryCount: retryCount, interval: interval}
}

// Connection should be implemented in handlers
func (rmqConfig *Config) Connection(url string) (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	for {
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Println("Failed to connect to RMQ: ", counts)
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = conn
			break
		}
		if counts > rmqConfig.retryCount {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func (rmqConfig *Config) CreateQueue(connection *amqp.Connection) (*amqp.Channel, error) {
	ch, err := connection.Channel()
	if err != nil {
		log.Panic("Failed to connect to RMQ: ", err)
		return nil, err
	}
	// defer ch.Close()
	// If channel closed after method execution, cannot publish messages
	exErr := ch.ExchangeDeclare(
		rmqConfig.exchange,
		"topic",
		true,  // durable?
		false, // auto-deleted?
		false, // internal?
		false, // no-wait?
		nil)
	if exErr != nil {
		log.Panic("RMQ Exchange Declare Error: ", exErr)
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		rmqConfig.qName, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Panic("RMQ Queue Declare Error: ", err)
		return nil, err
	}
	rmqBindError := ch.QueueBind(
		rmqConfig.qName,
		rmqConfig.rKey,
		rmqConfig.exchange,
		false,
		nil)
	if rmqBindError != nil {
		log.Panic("RMQ Bind Error: ", rmqBindError)
		return nil, err
	}
	log.Println("Queue Created: ", queue)
	return ch, nil
}

func (rmqConfig *Config) PublishMessage(channel *amqp.Channel, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), rmqConfig.interval*time.Second)
	defer cancel()
	err := channel.PublishWithContext(ctx,
		rmqConfig.exchange, // exchange
		rmqConfig.rKey,     // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Panic("RMQ Publish Error: ", err)
		return err
	}
	log.Printf(" [x] Sent %s\n", message)
	return nil
}

// ConsumeMessage This should be invoked in the main process
// initialize consumers for each queue, better to pass a an array
func (rmqConfig *Config) ConsumeMessage(ctx context.Context, channel *amqp.Channel, handler handler.UserHandler) error {
	msgs, err := channel.Consume(
		rmqConfig.qName, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		log.Panic("RMQ Consume Error: ", err)
		return err
	}
	var forever chan struct{}
	go func() {
		for d := range msgs {
			rmqMessage := RabbitmqMessage{delivery: d}
			// Invoke Function to Process Events
			log.Printf("Received a message: %s", d.Body)
			// HandleConsumer
			err = handler.HandleEvent(ctx, &rmqMessage)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	fmt.Println(<-forever)
	return nil
}

func (rmqConfig *Config) Reconnect(url string, conn *amqp.Connection) {
	go func() {
		for {
			// Check the connection status
			if conn.IsClosed() {
				log.Printf("Connection closed, reconnecting...")
				conn, err := rmqConfig.Connection(url)
				if err != nil {
					log.Printf("Failed to reconnect: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}
				log.Printf("Reconnected to RabbitMQ at %v", conn.LocalAddr())
			}

			// Send a ping to the server

			// Sleep for a while before checking again
			time.Sleep(30 * time.Second)
		}
	}()
}
