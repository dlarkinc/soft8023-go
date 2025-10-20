package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	const addr = "127.0.0.1:64001" // match the server

	// Connect
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	console := bufio.NewReader(os.Stdin)

	// Send the initial message (newline-terminated for the server's ReadString)
	if _, err := w.WriteString("This is from Client\n"); err != nil {
		log.Fatal(err)
	}
	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}

	for {
		// Read one line from server
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Server closed the connection.")
				return
			}
			log.Fatal("read error:", err)
		}
		fmt.Print("From Server: ", line)

		// Read user input (send as a line)
		out, _ := console.ReadString('\n')
		out = strings.TrimRight(out, "\r\n") // handle Windows \r\n and Unix \n

		// Send to server with newline
		if _, err := w.WriteString(out + "\n"); err != nil {
			log.Fatal("write error:", err)
		}
		if err := w.Flush(); err != nil {
			log.Fatal("flush error:", err)
		}

		// Exit condition
		if out == "bye" {
			break
		}
	}
}
