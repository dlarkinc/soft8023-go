package events

import "time"

type Event struct {
	Type   string    // e.g. "game.created", "game.played"
	GameID string    // e.g. "g-1"
	Msg    string    // free-form message (e.g. outcome text)
	Time   time.Time // filled by the bus on Publish
}
