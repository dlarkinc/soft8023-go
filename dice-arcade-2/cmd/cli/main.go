package main

import (
	"context"
	gamemanagerpb "dice-arcade/api/dicearcade/v1"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cli <game> (highlow|pig)")
		return
	}
	kind := os.Args[1]

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil { panic(err) }
	defer conn.Close()

	client := gamemanagerpb.NewGameManagerClient(conn)

	// Create
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cr, err := client.CreateGame(ctx, &gamemanagerpb.CreateGameRequest{Kind: kind})
	if err != nil { panic(err) }

	// Play once
	pr, err := client.PlayOnce(context.Background(), &gamemanagerpb.PlayOnceRequest{Id: cr.GetId()})
	if err != nil { panic(err) }

	fmt.Printf("[%s] %s\n", cr.GetId(), pr.GetOutcome())
}

