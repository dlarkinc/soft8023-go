package games

import (
	"dice-arcade/internal/dice"
	"fmt"
)

type HighLow struct{}

func (HighLow) Name() string { return "highlow" }

func (HighLow) PlayOnce() string {
	n := dice.D6()
	if n >= 4 {
		return fmt.Sprintf("HighLow: rolled %d → WIN", n)
	}
	return fmt.Sprintf("HighLow: rolled %d → LOSE", n)
}
