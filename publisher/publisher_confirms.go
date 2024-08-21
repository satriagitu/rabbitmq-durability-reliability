package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Confirm(false)
	failOnError(err, "Failed to put channel in confirm mode")

	ack, nack := ch.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

	body := "Hello, World!"
	err = ch.Publish(
		"exchange_name", // exchange
		"routing_key",   // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Set message as persistent
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	select {
	case <-ack:
		log.Println("Message confirmed")
	case <-nack:
		log.Println("Message not confirmed")
	}
}
