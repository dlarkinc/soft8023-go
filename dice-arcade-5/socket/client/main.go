package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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

	// Send getcount once, print the reply, exit.
	_, _ = fmt.Fprintln(conn, "getcount")
	reply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(reply)
}
