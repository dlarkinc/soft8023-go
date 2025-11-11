package decorate

import (
	"log"
)

type Game interface {
	Name() string
	PlayOnce() string
}

type Logging struct {
	Inner Game
}

func (d Logging) Name() string { return d.Inner.Name() }

func (d Logging) PlayOnce() string {
	out := d.Inner.PlayOnce()
	log.Printf("[LOG] %s → %s", d.Name(), out)
	return out
}

func WrapLogging(g Game) Game { return Logging{Inner: g} }
