package main

import (
	"fmt"
	"os"

	"dice-arcade/internal/manager"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cli <game> (highlow|pig)")
		return
	}
	kind := os.Args[1]

	m := manager.Get()
	id, g, err := m.Create(kind)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Play once and print outcome
	fmt.Printf("[%s] %s\n", id, g.PlayOnce())
}
