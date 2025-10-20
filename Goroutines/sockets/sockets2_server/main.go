package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync/atomic"
)

func main() {
	const addr = "127.0.0.1:64001"

	ln, err := net.Listen("tcp4", addr) // IPv4 for Windows friendliness
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Server started")
	fmt.Println("Waiting for client request..")

	var counter uint64

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		id := atomic.AddUint64(&counter, 1)

		fmt.Println("Connection no.", id)
		fmt.Println("New connection added:", conn.RemoteAddr())

		go handleConn(id, conn)
	}
}

func handleConn(id uint64, conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection from:", conn.RemoteAddr())

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		// Read one chunk/line; use ReadString to mirror simple text messages
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client at", conn.RemoteAddr(), "disconnected...")
			} else {
				log.Println("read error:", err)
			}
			return
		}

		msg := strings.TrimRight(line, "\r\n")
		if msg == "bye" {
			fmt.Println("Client at", conn.RemoteAddr(), "disconnected...")
			return
		}

		fmt.Println("from client", msg)

		// Echo back with newline (matches simple text protocol)
		if _, err := w.WriteString(msg + "\n"); err != nil {
			log.Println("write error:", err)
			return
		}
		if err := w.Flush(); err != nil {
			log.Println("flush error:", err)
			return
		}
	}
}
