package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the queue to ensure it exists
	q, err := ch.QueueDeclare(
		"my-queue", // queue name must match the sender
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnErr(err, "Failed to declare a queue")

	// Set QoS (one unacked message at a time)
	err = ch.Qos(1, 0, false)
	failOnErr(err, "Failed to set QoS")

	// Start consuming
	msgs, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer tag
		false,  // auto-ack = false (we’ll ack manually)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnErr(err, "Failed to register a consumer")

	log.Println("[*] Waiting for messages. To exit press CTRL+C")

	// Channel receive loop
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
		msg.Ack(false) // acknowledge successful processing
	}
}
