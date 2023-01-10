package amqp

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"time"
)

type Config struct {
	Rabbit *amqp.Connection
}

func Connection() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	for {
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq/")
		if err != nil {
			log.Println("Failed to connect to RMQ: ", counts)
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = conn
			break
		}
		if counts > 5 {
			fmt.Println(err)
			defer connection.Close()
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}

func CreateQueue(connection *amqp.Connection) (*amqp.Channel, error) {
	ch, err := connection.Channel()
	if err != nil {
		log.Panic("Failed to connect to RMQ: ", err)
		return nil, err
	}
	// defer ch.Close()
	// If channel closed after method execution, cannot publish messages
	exErr := ch.ExchangeDeclare(
		"rmq-ex",
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
		"rmq-q", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panic("RMQ Queue Declare Error: ", err)
		return nil, err
	}
	rmqBindError := ch.QueueBind(
		"rmq-q",
		"rmq-rk",
		"rmq-ex",
		false,
		nil)
	if rmqBindError != nil {
		log.Panic("RMQ Bind Error: ", rmqBindError)
		return nil, err
	}
	log.Println("Queue Created: ", queue)
	return ch, nil
}

func PublishMessage(channel *amqp.Channel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := "Hello World!"
	err := channel.PublishWithContext(ctx,
		"rmq-ex", // exchange
		"rmq-rk", // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Panic("RMQ Publish Error: ", err)
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}

func ConsumeMessage(channel *amqp.Channel) error {
	msgs, err := channel.Consume(
		"rmq-q", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Panic("RMQ Consume Error: ", err)
		return err
	}
	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	fmt.Println(<-forever)
	return nil
}

func Reconnect(conn *amqp.Connection) {
	go func() {
		for {
			// Check the connection status
			if conn.IsClosed() {
				log.Printf("Connection closed, reconnecting...")
				conn, err := Connection()
				if err != nil {
					log.Printf("Failed to reconnect: %v", err)
					time.Sleep(5 * time.Second)
					continue
				}
				log.Printf("Reconnected to RabbitMQ at %v", conn.LocalAddr())
			}

			// Send a ping to the server
			err := conn.
			if err != nil {
				log.Printf("Ping failed: %v", err)
				conn.Close()
			}

			// Sleep for a while before checking again
			time.Sleep(30 * time.Second)
		}
	}()
}
