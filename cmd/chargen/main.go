package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/script-wizards/chargen/pkg/character"
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	r := chi.NewRouter()

	r.Get("/", handler)

	log.Fatal(http.ListenAndServe(":"+port, r))
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

func calcInitiative(score int) int {
	switch score {
	case 3:
		return -2
	case 4, 5, 6, 7, 8:
		return -1
	case 9, 10, 11, 12:
		return 0
	case 13, 14, 15, 16, 17:
		return 1
	case 18:
		return 2
	default:
		return 0
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
		return "illiterate"
	case score >= 6 && score <= 18:
		return "literate"
	default:
		return ""
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := os.ReadFile("templates/template.html")
	check(err)

	t, err := template.New("webpage").Parse(string(tpl))
	check(err)

	char := character.NewRandomChar()

	data := Character{
		Class:        char.Class,
		Level:        char.Level,
		Title:        char.Title,
		Alignment:    "Neutral",
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
		ArmorClass:   5,
		Initiative:   calcInitiative(char.DEX),
		SaveDeath:    char.SaveDeath,
		SaveWands:    char.SaveWands,
		SaveParalyze: char.SaveParalyze,
		SaveBreath:   char.SaveBreath,
		SaveSpells:   char.SaveSpells,
		OpenDoors:    calcOpenDoors(char.STR),
		THAC0:        calcTHAC0(19),
		XPBonus:      character.PrimeRequisite(char.Class, char.STR, char.INT, char.WIS, char.DEX, char.CON, char.CHA),
		XPNext:       char.NextLevel(),
		Inventory:    char.InventoryString(),
	}

	err = t.Execute(w, data)
	check(err)
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
	Inventory    string
}
