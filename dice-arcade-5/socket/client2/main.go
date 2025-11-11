package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	addr := "localhost:9000"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)

	// Background reader: prints any lines from the server (greeting + replies).
	go func() {
		for {
			line, err := serverReader.ReadString('\n')
			if err != nil {
				// Server closed or network error; exit the process.
				fmt.Println("\nDisconnected from server.")
				os.Exit(0)
			}
			fmt.Print(strings.TrimRight(line, "\n") + "\n")
			// Reprint prompt after async server message
			fmt.Print("> ")
		}
	}()

	// Simple REPL for user commands
	input := bufio.NewScanner(os.Stdin)
	fmt.Println("Type 'getcount' to query, or 'quit' to exit.")
	fmt.Print("> ")

	for input.Scan() {
		cmd := strings.TrimSpace(input.Text())
		if cmd == "" {
			fmt.Print("> ")
			continue
		}
		if strings.EqualFold(cmd, "quit") {
			return
		}

		// Send the command to server
		if _, err := fmt.Fprintln(serverWriter, cmd); err != nil {
			log.Fatalf("write: %v", err)
		}
		if err := serverWriter.Flush(); err != nil {
			log.Fatalf("flush: %v", err)
		}

		// prompt again; the background reader will print the server's reply
		fmt.Print("> ")
	}

	if err := input.Err(); err != nil {
		log.Fatalf("stdin: %v", err)
	}
}
