package games

import "fmt"

func New(kind string) (Game, error) {
	switch kind {
	case "highlow":
		return HighLow{}, nil
	case "pig":
		return Pig{}, nil
	default:
		return nil, fmt.Errorf("unknown game: %s", kind)
	}
}
