package character

import (
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Character struct {
	Class     string
	Level     int
	Title     string
	Alignment string
	STR       int
	INT       int
	WIS       int
	DEX       int
	CON       int
	CHA       int
}

func roll(die int) int {
	return rng.Intn(die-1) + 1
}

func roll3d6() int {
	return roll(6) + roll(6) + roll(6)
}

func NewRandomChar() *Character {
	return &Character{
		STR: roll3d6(),
		INT: roll3d6(),
		WIS: roll3d6(),
		DEX: roll3d6(),
		CON: roll3d6(),
		CHA: roll3d6(),
	}
}
