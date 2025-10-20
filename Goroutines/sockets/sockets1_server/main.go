package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	const (
		HOST = "127.0.0.1"
		PORT = "64000"
	)

	// Listen for incoming connections on TCP port 64000
	ln, err := net.Listen("tcp", net.JoinHostPort(HOST, PORT))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("Server listening on %s:%s\n", HOST, PORT)

	// Accept a single client connection
	conn, err := ln.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected from:", conn.RemoteAddr())

	// Echo loop
	buffer := make([]byte, 512)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}
			panic(err)
		}
		if n == 0 {
			break
		}
		// Echo back to the client
		_, err = conn.Write(buffer[:n])
		if err != nil {
			panic(err)
		}
	}
}
