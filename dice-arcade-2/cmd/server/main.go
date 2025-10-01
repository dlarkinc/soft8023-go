package main

import (
	"context"
	gamemanagerpb "dice-arcade/api/dicearcade/v1"
	"dice-arcade/internal/manager"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	gamemanagerpb.UnimplementedGameManagerServer
	mgr manager.Manager
}

func (s *server) CreateGame(ctx context.Context, req *gamemanagerpb.CreateGameRequest) (*gamemanagerpb.CreateGameResponse, error) {
	id, g, err := s.mgr.Create(req.GetKind())
	if err != nil {
		return nil, err
	}
	return &gamemanagerpb.CreateGameResponse{Id: id, Name: g.Name()}, nil
}

func (s *server) PlayOnce(ctx context.Context, req *gamemanagerpb.PlayOnceRequest) (*gamemanagerpb.PlayOnceResponse, error) {
	g, ok := s.mgr.Get(req.GetId())
	if !ok {
		return nil, fmt.Errorf("game not found")
	}
	return &gamemanagerpb.PlayOnceResponse{Outcome: g.PlayOnce()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := &server{mgr: manager.Get()}
	grpcServer := grpc.NewServer()
	gamemanagerpb.RegisterGameManagerServer(grpcServer, s)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
