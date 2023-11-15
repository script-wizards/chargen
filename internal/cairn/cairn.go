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
	Details    []string
	Gear       []string
	Questions  []string
	Answers    []string
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

func HandleBlank(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/cairn-blank.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, nil); err != nil {
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

func bonds() []string {
	lines := make([]string, 5)
	src := splitLines(bondsList[rng.Intn(len(bondsList))], 57)
	copy(lines, src)
	return lines
}

func omens() []string {
	lines := make([]string, 5)
	src := splitLines(omensList[rng.Intn(len(omensList))], 57)
	copy(lines, src)
	return lines
}

func details(background string) []string {
	lines := make([]string, 5)
	src := splitLines(backgroundDetails[background], 57)
	copy(lines, src)
	return lines
}

func gear(background string) []string {
	return backgroundGear[background]
}

func questions(background string) []string {
	return backgroundQuestions[background]
}

func answers(background string) []string {
	answers := make([]string, 2)
	for i, answer := range backgroundAnswers[background] {
		answers[i] = answer[rng.Intn(len(answer))]
	}
	return answers
}

func NewCairnCharacter() *Character {
	background := background()
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
		Bonds:      bonds(),
		Omens:      omens(),
		Gold:       dice.Roll3d6(),
		Details:    details(background),
		Gear:       gear(background),
		Questions:  questions(background),
		Answers:    answers(background),
	}
}
