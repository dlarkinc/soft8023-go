// cmd/cli/main.go
package main

import (
	"context"
	gamemanagerpb "dice-arcade/api/dicearcade/v1"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	player := flag.String("player", "", "player id (e.g. alice)")
	addr := flag.String("addr", "localhost:50051", "server address")
	flag.Parse()

	if *player == "" {
		log.Fatal("missing --player")
	}
	if flag.NArg() < 1 {
		log.Fatalf("usage: cli --player=<id> <create|join|play|state> [args]")
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := gamemanagerpb.NewGameManagerClient(conn)

	cmd := flag.Arg(0)
	switch cmd {

	case "create":
		if flag.NArg() < 2 {
			log.Fatal("usage: ... create <kind>")
		}
		kind := flag.Arg(1)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cr, err := c.CreateGame(ctx, &gamemanagerpb.CreateGameRequest{Kind: kind})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("created game: %s (%s)\n", cr.GetId(), cr.GetName())

		// auto-join creator (optional)
		_, err = c.JoinGame(context.Background(), &gamemanagerpb.JoinGameRequest{
			Id: cr.GetId(), PlayerId: *player,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("joined as:", *player)

	case "join":
		if flag.NArg() < 2 {
			log.Fatal("usage: ... join <gameId>")
		}
		id := flag.Arg(1)
		jr, err := c.JoinGame(context.Background(), &gamemanagerpb.JoinGameRequest{
			Id: id, PlayerId: *player,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("players: %v (you are #%d)\n", jr.GetPlayers(), jr.GetYourIndex())

	case "play":
		if flag.NArg() < 2 {
			log.Fatal("usage: ... play <gameId>")
		}
		id := flag.Arg(1)
		pr, err := c.PlayTurn(context.Background(), &gamemanagerpb.PlayTurnRequest{
			Id: id, PlayerId: *player,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("outcome: %s\nnext: %s\n", pr.GetOutcome(), pr.GetNextPlayerId())

	case "state":
		if flag.NArg() < 2 {
			log.Fatal("usage: ... state <gameId>")
		}
		id := flag.Arg(1)
		sr, err := c.GetGameState(context.Background(), &gamemanagerpb.GetGameStateRequest{Id: id})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("game: %s\nplayers: %v\ncurrent: %s\n",
			sr.GetGameName(), sr.GetPlayers(), sr.GetCurrentPlayerId())

	default:
		log.Fatalf("unknown command: %s", cmd)
	}
}
