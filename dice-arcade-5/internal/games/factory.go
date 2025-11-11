package games

import (
	"dice-arcade/internal/games/decorate"
	"dice-arcade/internal/games/legacy"
	"fmt"
)

func New(kind string) (Game, error) {
	switch kind {
	case "highlow":
		return HighLow{}, nil
	case "pig":
		return Pig{}, nil
		// internal/games/factory.go (add case)
	case "highlow_legacy":
		l := legacy.NewHighLowLegacy()
		return legacy.NewHighLowAdapter(l), nil
	case "highlow+log":
		g := HighLow{}
		return decorate.WrapLogging(g), nil
	case "pig+score":
		g := Pig{}
		return decorate.WrapScoring(g), nil
	case "highlow+log+score":
		g := HighLow{}
		return decorate.WrapLogging(decorate.WrapScoring(g)), nil
	default:
		return nil, fmt.Errorf("unknown game: %s", kind)
	}
}
