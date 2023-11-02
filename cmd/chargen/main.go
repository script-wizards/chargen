package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi"
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

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := os.ReadFile("template.html")
	check(err)

	t, err := template.New("webpage").Parse(string(tpl))
	check(err)

	data := Character{
		Class:        "Halfling",
		Level:        1,
		Title:        "Halfling Swordmaster",
		Alignment:    "Neutral",
		STR:          12,
		INT:          9,
		WIS:          13,
		DEX:          16,
		CON:          14,
		CHA:          15,
		ModSTR:       calcMod(12),
		ModINT:       calcINTMod(9),
		ModWIS:       calcMod(13),
		ModDEX:       calcMod(16),
		ModCON:       calcMod(14),
		ModCHA:       calcCHAMod(15),
		Literacy:     "literate",
		HitPoints:    6,
		HitDie:       6,
		ArmorClass:   5,
		Initiative:   calcInitiative(16),
		SaveDeath:    8,
		SaveWands:    9,
		SaveParalyze: 10,
		SaveBreath:   13,
		SaveSpells:   12,
		OpenDoors:    calcOpenDoors(12),
		THAC0:        calcTHAC0(19),
		XPBonus:      5,
		XPNext:       2000,
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
}
