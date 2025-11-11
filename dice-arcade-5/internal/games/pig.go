package games

import (
	"dice-arcade/internal/dice"
	"fmt"
)

type Pig struct{}

func (Pig) Name() string { return "pig" }

func (Pig) PlayOnce() string {
	n := dice.D6()
	if n == 1 {
		return fmt.Sprintf("Pig: rolled %d → BUST", n)
	}
	return fmt.Sprintf("Pig: rolled %d → +points", n)
}
