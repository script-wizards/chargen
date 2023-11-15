package dice

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func Roll(die int) int {
	return 1 + rng.Intn(die)
}

func Roll3d6() int {
	return Roll(6) + Roll(6) + Roll(6)
}
