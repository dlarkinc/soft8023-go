package main

import (
	"fmt"
)

func ping(messages chan string) {
	for i := 0; i < 5; i++ {
		messages <- "ping"
	}
	close(messages) // close channel when done sending
}

func main() {
	messages := make(chan string)

	go ping(messages)

	for msg := range messages {
		fmt.Println(msg)
	}
}
