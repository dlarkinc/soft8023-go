// cmd/cli/main.go
package main

import (
	"context"
	gamemanagerpb "dice-arcade/api/dicearcade/v1"
	"dice-arcade/internal/manager/proxy" // ← add this import
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
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Wrap the generated client with the proxy (e.g., 5 calls/sec throttle)
	raw := gamemanagerpb.NewGameManagerClient(conn)
	client := proxy.NewClientProxy(raw, 5)

	// Create
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cr, err := client.CreateGame(ctx, &gamemanagerpb.CreateGameRequest{Kind: kind})
	if err != nil {
		panic(err)
	}

	// Play once
	pr, err := client.PlayOnce(context.Background(), &gamemanagerpb.PlayOnceRequest{Id: cr.GetId()})
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] %s\n", cr.GetId(), pr.GetOutcome())

	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()

	sum, err := client.GetSummary(ctx2, &gamemanagerpb.GetSummaryRequest{Id: cr.GetId()})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Summary for %s (%s):\n", sum.GetGameId(), sum.GetGameName())
	for _, r := range sum.GetRolls() {
		fmt.Printf("  Roll %d: %d → %s\n", r.GetRollNumber(), r.GetValue(), r.GetOutcome())
	}
	fmt.Printf("Total rolls: %d, Wins: %d\n", sum.GetTotalRolls(), sum.GetTotalWins())

}
