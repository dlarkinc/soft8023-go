package main

import (
	"fmt"
	"net"
)

func main() {
	const (
		HOST = "localhost"
		PORT = "64000"
	)

	conn, err := net.Dial("tcp", net.JoinHostPort(HOST, PORT))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send data
	message := []byte("Echo!!!!")
	_, err = conn.Write(message)
	if err != nil {
		panic(err)
	}

	// Receive response (up to 512 bytes)
	buffer := make([]byte, 512) // buffer "slice" - allocating 512 bytes of memory
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received:", string(buffer[:n]))
}
