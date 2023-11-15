package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"text/template"
	"time"

	"github.com/go-chi/chi"
	"github.com/script-wizards/chargen/internal/cairn"
	"github.com/script-wizards/chargen/internal/character"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	r := chi.NewRouter()

	r.Get("/", handler)
	r.Get("/blank", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/blank.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	r.Get("/cairn", cairn.Handler)
	r.Get("/cairn-blank", cairn.HandleBlank)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/404.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := struct {
			Msg string
		}{
			Msg: random404(),
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func random404() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	messages := []string{
		"Impossible. Perhaps the archives are incomplete.",
		"This page is too strong for you, traveler.",
		"This page is in another castle.",
		"One does not simply walk into this page.",
		"It's dangerous to go alone, take this link back home.",
		"\"DID YOU PUT THIS PAGE IN THE GOBLET OF FIRE\", Dumbledore said calmly.",
	}
	return messages[rng.Intn(len(messages))]
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func calcMod(score int) int {
	switch score {
	case 3:
		return -3
	case 4, 5:
		return -2
	case 6, 7, 8:
		return -1
	case 9, 10, 11, 12:
		return 0
	case 13, 14, 15:
		return 1
	case 16, 17:
		return 2
	case 18:
		return 3
	default:
		return 0
	}
}

func calcINTMod(score int) int {
	switch score {
	case 3, 4, 5, 6, 7, 8, 9, 10, 11, 12:
		return 0
	case 13, 14, 15:
		return 1
	case 16, 17:
		return 2
	case 18:
		return 3
	default:
		return 0
	}
}

func calcCHAMod(score int) []int {
	switch score {
	case 3:
		return []int{-2, 1, 4}
	case 4, 5:
		return []int{-1, 2, 5}
	case 6, 7, 8:
		return []int{-1, 3, 6}
	case 9, 10, 11, 12:
		return []int{0, 4, 7}
	case 13, 14, 15:
		return []int{1, 5, 8}
	case 16, 17:
		return []int{1, 6, 9}
	case 18:
		return []int{1, 7, 10}
	default:
		return nil
	}
}

func calcOpenDoors(score int) int {
	switch score {
	case 3, 4, 5, 6, 7, 8:
		return 1
	case 9, 10, 11, 12:
		return 2
	case 13, 14, 15:
		return 3
	case 16, 17:
		return 4
	case 18:
		return 5
	default:
		return 0
	}
}

func calcTHAC0(score int) string {
	vals := make([]int, 0, 10)
	for i := 0; i < 10; i++ {
		vals = append(vals, score-(9-i))
	}
	str := ""
	for _, v := range vals {
		str += fmt.Sprintf("%-2d", v) + " "
	}
	return str
}

func literacy(score int) string {
	switch {
	case score >= 3 && score <= 5:
		return "Illiterate"
	case score >= 6 && score <= 18:
		return "Literate"
	default:
		return ""
	}
}

var validClasses = []string{
	"cleric",
	"dwarf",
	"elf",
	"fighter",
	"halfling",
	"magic-user",
	"thief",
}

func handler(w http.ResponseWriter, r *http.Request) {
	class := ""
	if r.URL.Query().Get("class") != "" {
		class = strings.ToLower(r.URL.Query().Get("class"))
		if !slices.Contains(validClasses, class) {
			class = ""
		}
	}

	tpl, err := os.ReadFile("templates/template.html")
	check(err)

	t, err := template.New("webpage").Parse(string(tpl))
	check(err)

	var char *character.Character
	if class != "" {
		char = character.NewCharClass(class)
	} else {
		char = character.NewRandomChar()
	}

	caser := cases.Title(language.English)

	data := Character{
		Class:        caser.String(char.Class),
		Level:        char.Level,
		Title:        char.Title,
		Alignment:    char.Alignment,
		STR:          char.STR,
		INT:          char.INT,
		WIS:          char.WIS,
		DEX:          char.DEX,
		CON:          char.CON,
		CHA:          char.CHA,
		ModSTR:       calcMod(char.STR),
		ModINT:       calcINTMod(char.INT),
		ModWIS:       calcMod(char.WIS),
		ModDEX:       calcMod(char.DEX),
		ModCON:       calcMod(char.CON),
		ModCHA:       calcCHAMod(char.CHA),
		Literacy:     literacy(char.INT),
		HitPoints:    char.HitPoints,
		HitDie:       char.HitDie,
		ArmorClass:   char.SetAC(),
		Initiative:   char.Initiative(),
		SaveDeath:    char.SaveDeath,
		SaveWands:    char.SaveWands,
		SaveParalyze: char.SaveParalyze,
		SaveBreath:   char.SaveBreath,
		SaveSpells:   char.SaveSpells,
		OpenDoors:    calcOpenDoors(char.STR),
		THAC0:        calcTHAC0(19),
		XPBonus:      character.PrimeRequisite(char.Class, char.STR, char.INT, char.WIS, char.DEX, char.CON, char.CHA),
		XPNext:       char.NextLevel(),
		Gold:         char.Gold,
		Inventory:    padString(char.Inventory),
		Abilities:    padString(char.Abilities),
	}

	err = t.Execute(w, data)
	check(err)
}

func padString(list []string) string {
	s := ""
	for i, item := range list {
		if i == len(list)-1 {
			s += " " + item
		} else {
			s += " " + item + "\n"
		}
	}
	return s
}

type Character struct {
	Class        string
	Level        int
	Title        string
	Alignment    string
	STR          int
	INT          int
	WIS          int
	DEX          int
	CON          int
	CHA          int
	ModSTR       int
	ModINT       int
	ModWIS       int
	ModDEX       int
	ModCON       int
	ModCHA       []int
	Literacy     string
	HitPoints    int
	HitDie       int
	ArmorClass   int
	Initiative   int
	SaveDeath    int
	SaveWands    int
	SaveParalyze int
	SaveBreath   int
	SaveSpells   int
	OpenDoors    int
	THAC0        string
	XPBonus      int
	XPNext       int
	Gold         int
	Inventory    string
	Abilities    string
}
