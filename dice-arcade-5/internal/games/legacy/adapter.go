package legacy

import "fmt"

// Target interface lives in internal/games (we duplicate minimal shape to avoid import cycle)
// In your code, import the real games package and implement that interface.
type Game interface {
	Name() string
	PlayOnce() string
}

type HighLowAdapter struct {
	legacy *HighLowLegacy
}

func NewHighLowAdapter(l *HighLowLegacy) *HighLowAdapter {
	return &HighLowAdapter{legacy: l}
}

func (a *HighLowAdapter) Name() string { return "highlow_legacy" }

func (a *HighLowAdapter) PlayOnce() string {
	roll, win := a.legacy.RunOnce()
	if win {
		return fmt.Sprintf("HighLow(legacy): rolled %d → WIN", roll)
	}
	return fmt.Sprintf("HighLow(legacy): rolled %d → LOSE", roll)
}
