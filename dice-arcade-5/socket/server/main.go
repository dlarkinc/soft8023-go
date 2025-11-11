package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/streadway/amqp"
)

const (
	amqpURL = "amqp://admin:password@localhost:5672/"
	tcpAddr     = ":9000"
)

var (
	mu       sync.Mutex
	msgCount uint64
)

func main() {

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("amqp dial: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("amqp channel: %v", err)
	}
	defer ch.Close()

	// Ensure exchange exists (topic)
	if err := ch.ExchangeDeclare("game.events", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("exchange declare: %v", err)
	}

	// Ephemeral queue for this server
	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		log.Fatalf("queue declare: %v", err)
	}

	// Bind to common event patterns
	if err := ch.QueueBind(q.Name, "player.*", "game.events", false, nil); err != nil {
		log.Fatalf("queue bind player.*: %v", err)
	}
	if err := ch.QueueBind(q.Name, "game.*", "game.events", false, nil); err != nil {
		log.Fatalf("queue bind game.*: %v", err)
	}

	deliveries, err := ch.Consume(q.Name, "", true, true, false, false, nil)
	if err != nil {
		log.Fatalf("consume: %v", err)
	}

	// --- TCP server ---
	ln, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Socket server listening on %s", tcpAddr)

	// Consume messages in background
	go func() {
		for d := range deliveries {
			payload := strings.TrimSpace(string(d.Body))
			line := fmt.Sprintf("[%s] %s", d.RoutingKey, payload)
			log.Println(line)

			mu.Lock()
			msgCount++
			mu.Unlock()
		}
		log.Println("RabbitMQ consumer ended")
	}()

	// Accept multiple clients; each can issue "getcount"
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return // client disconnected
		}
		cmd := strings.TrimSpace(line)
		if strings.EqualFold(cmd, "getcount") {
			mu.Lock()
			count := msgCount
			mu.Unlock()
			_, _ = fmt.Fprintf(w, "Total number of messages so far: %d\n", count)
			_ = w.Flush()
			// keep the connection open so the client can ask again if they want
		} else if cmd == "" {
			continue
		} else {
			_, _ = w.WriteString("Unrecognized command. Try: getcount\n")
			_ = w.Flush()
		}
	}
}
