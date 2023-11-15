package cairn

import (
	"math/rand"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/script-wizards/chargen/internal/dice"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Character struct {
	Name       string
	Background string
	STR        int
	DEX        int
	WIL        int
	HP         int
	Armor      int
	Traits     []string
	Age        int
	Bonds      []string
	Omens      []string
	Inventory  []string
	Gold       int
}

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/cairn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := NewCairnCharacter()
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// splitLines splits a string into lines of a given length, breaking on spaces
func splitLines(str string, linelen int) []string {
	var lines []string
	var line string
	words := strings.Split(str, " ")
	for _, word := range words {
		if len(line)+len(word) > linelen {
			lines = append(lines, line)
			line = ""
		}
		line += word + " "
	}
	lines = append(lines, line)
	return lines
}

func name(background string) string {
	namesList := names[background]
	return namesList[rng.Intn(len(namesList))]
}

func background() string {
	return backgrounds[rng.Intn(len(backgrounds))]
}

func trait() []string {
	caser := cases.Title(language.English)
	traitsList := make([]string, len(traits))
	for i, trait := range traits {
		traitsList[i] = caser.String(trait[rng.Intn(len(trait))])
	}
	return traitsList
}

func NewCairnCharacter() *Character {
	background := background()
	bond := "You consumed a mischievous spirit that periodically wreaks havoc on your insides, demanding to be taken home. It wants to keep you alive, at least until it is free. It can detect magic and knows quite a bit about The Woods."
	bondLines := splitLines(bond, 60)
	omen := "It feels like winter has arrived too quickly this year, frost and snows making their appearance much earlier than expected. There is talk of a pattern to the frost found in windows, ponds, and cracks in the ground. It almost looks like a map."
	omenLines := splitLines(omen, 60)
	caser := cases.Title(language.English)
	return &Character{
		Name:       name(background),
		Background: caser.String(background),
		STR:        dice.Roll3d6(),
		DEX:        dice.Roll3d6(),
		WIL:        dice.Roll3d6(),
		HP:         dice.Roll(6),
		Armor:      0,
		Traits:     trait(),
		Age:        dice.Roll(20) + dice.Roll(20) + 10,
		Bonds:      bondLines,
		Omens:      omenLines,
		Gold:       dice.Roll3d6(),
	}
}
