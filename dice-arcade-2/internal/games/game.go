package games

type Game interface {
	Name() string
	PlayOnce() string // returns outcome text (no I/O yet)
}
