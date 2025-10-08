package decorate

import "strings"

type Scoring struct {
	Inner  Game
	Points int
}

func (d *Scoring) Name() string { return d.Inner.Name() }

func (d *Scoring) PlayOnce() string {
	out := d.Inner.PlayOnce()
	if strings.Contains(out, "WIN") || strings.Contains(out, "+points") {
		d.Points += 1
	}
	return out
}

func WrapScoring(g Game) *Scoring { return &Scoring{Inner: g} }
