package cairn

import (
	"net/http"
	"text/template"
)

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
	Bonds      string
	Omens      string
	Inventory  []string
	Gold       int
}

func Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/cairn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := Character{
		Name:       "Landsford",
		Background: "Barber-Surgeon",
		STR:        13,
		DEX:        14,
		WIL:        9,
		HP:         6,
		Armor:      1,
		Traits: []string{
			"Statuesque",
			"Birthmarked",
			"Luxurious",
			"Rat-like",
			"Stuttering",
			"Elegant",
			"Disciplined",
			"Aggressive",
		},
		Age: 50,
		Bonds: `You consumed a mischievous spirit that periodically wreaks havoc on
		your insides, demanding to be taken home. It wants to keep you alive, at
		least until it is free. It can detect magic and knows quite a bit about The
		Woods.`,
		Omens: `It feels like winter has arrived too quickly this year, frost and snows
		making their appearance much earlier than expected. There is talk of a
		pattern to the frost found in windows, ponds, and cracks in the ground.
		It almost looks like a map.`,
		Inventory: []string{
			"Torch",
			"Rations",
			"16gp",
		},
		Gold: 15,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
