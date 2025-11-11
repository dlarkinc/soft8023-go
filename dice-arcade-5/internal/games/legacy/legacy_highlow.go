package legacy

import (
	"dice-arcade/internal/dice"
)

type HighLowLegacy struct {
	// legacy needs an injected roller, different method names, etc.
	Roll func() int
}

func NewHighLowLegacy() *HighLowLegacy {
	return &HighLowLegacy{Roll: dice.D6}
}

// Legacy API: returns (roll, win)
func (h *HighLowLegacy) RunOnce() (int, bool) {
	n := h.Roll()
	return n, n >= 4
}
