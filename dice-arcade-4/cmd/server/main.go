package main

import (
	"context"
	gamemanagerpb "dice-arcade/api/dicearcade/v1"
	"dice-arcade/internal/games"
	"dice-arcade/internal/manager"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type session struct {
	game    games.Game
	players []string
	turn    int // index into players
}

type server struct {
	gamemanagerpb.UnimplementedGameManagerServer
	mgr manager.Manager
	// sessions keyed by game id; protect with a mutex for simplicity
	mu       sync.Mutex
	sessions map[string]*session
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

func (s *server) GetSummary(ctx context.Context, req *gamemanagerpb.GetSummaryRequest) (*gamemanagerpb.GameSummaryResponse, error) {
	// Fake data for demo purposes
	summary := &gamemanagerpb.GameSummaryResponse{
		GameId:   req.GetId(),
		GameName: "highlow",
		Rolls: []*gamemanagerpb.GameSummaryResponse_RollResult{
			{RollNumber: 1, Value: 5, Outcome: gamemanagerpb.Outcome_OUTCOME_WIN},
			{RollNumber: 2, Value: 3, Outcome: gamemanagerpb.Outcome_OUTCOME_LOSE},
			{RollNumber: 3, Value: 6, Outcome: gamemanagerpb.Outcome_OUTCOME_WIN},
		},
		TotalRolls: 3,
		TotalWins:  2,
	}
	return summary, nil
}

func (s *server) JoinGame(ctx context.Context, req *gamemanagerpb.JoinGameRequest) (*gamemanagerpb.JoinGameResponse, error) {
	id, pid := req.GetId(), req.GetPlayerId()

	s.mu.Lock()
	defer s.mu.Unlock()

	sess := s.sessions[id]
	if sess == nil {
		// lazy create if game exists but no session (a real, just in case)
		if g, ok := s.mgr.Get(id); ok {
			sess = &session{game: g}
			s.sessions[id] = sess
		} else {
			return nil, fmt.Errorf("game not found")
		}
	}
	// add if missing
	idx := -1
	for i, p := range sess.players {
		if p == pid {
			idx = i
			break
		}
	}
	if idx == -1 {
		sess.players = append(sess.players, pid)
		idx = len(sess.players) - 1
	}
	return &gamemanagerpb.JoinGameResponse{
		Players:   sess.players,
		YourIndex: int32(idx),
	}, nil
}

func (s *server) GetGameState(ctx context.Context, req *gamemanagerpb.GetGameStateRequest) (*gamemanagerpb.GetGameStateResponse, error) {
	id := req.GetId()
	s.mu.Lock()
	sess := s.sessions[id]
	s.mu.Unlock()
	if sess == nil {
		return nil, fmt.Errorf("game not found")
	}

	current := ""
	if len(sess.players) > 0 {
		current = sess.players[sess.turn%len(sess.players)]
	}
	return &gamemanagerpb.GetGameStateResponse{
		Players:         append([]string(nil), sess.players...),
		CurrentPlayerId: current,
		GameName:        sess.game.Name(),
	}, nil
}

func (s *server) PlayTurn(ctx context.Context, req *gamemanagerpb.PlayTurnRequest) (*gamemanagerpb.PlayTurnResponse, error) {
	id, pid := req.GetId(), req.GetPlayerId()

	s.mu.Lock()
	defer s.mu.Unlock()

	sess := s.sessions[id]
	if sess == nil {
		return nil, fmt.Errorf("game not found")
	}
	if len(sess.players) == 0 {
		return nil, fmt.Errorf("no players joined")
	}

	current := sess.players[sess.turn%len(sess.players)]
	if pid != current {
		return nil, fmt.Errorf("not your turn; current: %s", current)
	}

	// play once
	outcome := sess.game.PlayOnce()

	// advance turn
	sess.turn++
	next := sess.players[sess.turn%len(sess.players)]

	return &gamemanagerpb.PlayTurnResponse{
		Outcome:      outcome,
		NextPlayerId: next,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	s := &server{mgr: manager.Get(), sessions: make(map[string]*session)}
	grpcServer := grpc.NewServer()
	gamemanagerpb.RegisterGameManagerServer(grpcServer, s)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
