package proxy

import (
	"context"
	gamemanagerpb "dice-arcade/api/dicearcade/v1"
	"time"
)

type ClientProxy struct {
	Inner     gamemanagerpb.GameManagerClient
	MaxPerSec int // simple token-ish limiter
	lastCall  time.Time
	minGap    time.Duration
}

func NewClientProxy(inner gamemanagerpb.GameManagerClient, maxPerSec int) *ClientProxy {
	gap := time.Second
	if maxPerSec > 0 {
		gap = time.Second / time.Duration(maxPerSec)
	}
	return &ClientProxy{Inner: inner, MaxPerSec: maxPerSec, minGap: gap}
}

func (p *ClientProxy) throttle() {
	now := time.Now()
	wait := p.minGap - now.Sub(p.lastCall)
	if wait > 0 {
		time.Sleep(wait)
	}
	p.lastCall = time.Now()
}

// Pass-through methods add proxy behavior:
func (p *ClientProxy) CreateGame(ctx context.Context, req *gamemanagerpb.CreateGameRequest, opts ...interface{}) (*gamemanagerpb.CreateGameResponse, error) {
	p.throttle()
	// NB: opts isn’t the exact grpc.CallOption type here to keep snippet light;
	// in your code, import "google.golang.org/grpc" and use ...grpc.CallOption.
	return p.Inner.CreateGame(ctx, req /* opts... */)
}

func (p *ClientProxy) PlayOnce(ctx context.Context, req *gamemanagerpb.PlayOnceRequest, opts ...interface{}) (*gamemanagerpb.PlayOnceResponse, error) {
	p.throttle()
	return p.Inner.PlayOnce(ctx, req /* opts... */)
}

func (p *ClientProxy) GetSummary(ctx context.Context, req *gamemanagerpb.GetSummaryRequest, opts ...interface{}) (*gamemanagerpb.GameSummaryResponse, error) {
	p.throttle()
	return p.Inner.GetSummary(ctx, req /* opts... */)
}
