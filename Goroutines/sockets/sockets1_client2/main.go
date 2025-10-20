package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	const (
		HOST = "127.0.0.1" // use IPv4 to avoid Windows IPv6 issues
		PORT = "64000"
	)

	conn, err := net.Dial("tcp4", net.JoinHostPort(HOST, PORT))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	console := bufio.NewReader(os.Stdin)

	response := "Y"
	dataToSend := "Echo!!!!"

	for strings.EqualFold(response, "Y") {
		// Send message to server
		if _, err := conn.Write([]byte(dataToSend)); err != nil {
			if err == io.EOF {
				fmt.Println("Server closed connection")
				return
			}
			panic(err)
		}

		// Receive up to 512 bytes
		buffer := make([]byte, 512)
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed connection")
				return
			}
			panic(err)
		}
		fmt.Println("Received:", string(buffer[:n]))

		// Ask user if they want to send again
		fmt.Print("Send again (Y/N)? ")
		respLine, _ := console.ReadString('\n')
		response = strings.TrimSpace(respLine)

		// If yes, prompt for new message
		if strings.EqualFold(response, "Y") {
			fmt.Print("Enter data: ")
			line, _ := console.ReadString('\n')
			dataToSend = strings.TrimRight(line, "\r\n") // trims both Unix and Windows newlines
		}
	}

	fmt.Println("Goodbye!")
}
