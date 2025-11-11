package dice

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func D6() int { return rng.Intn(6) + 1 }
