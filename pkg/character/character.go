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
	Inventory []string
}

func roll(die int) int {
	return rng.Intn(die-1) + 1
}

func roll3d6() int {
	return roll(6) + roll(6) + roll(6)
}

func NewRandomChar() *Character {
	// roll stats first before picking class
	STR := roll3d6()
	INT := roll3d6()
	WIS := roll3d6()
	DEX := roll3d6()
	CON := roll3d6()
	CHA := roll3d6()

	class := pickClass(STR, INT, WIS, DEX, CON, CHA)
	return &Character{
		Class:     class,
		Level:     1,
		Title:     generateTitle(class, 1),
		STR:       STR,
		INT:       INT,
		WIS:       WIS,
		DEX:       DEX,
		CON:       CON,
		CHA:       CHA,
		Inventory: generateInventory(class),
	}
}

func (c *Character) InventoryString() string {
	s := ""
	for _, item := range c.Inventory {
		s += "  " + item + "\n"
	}
	return s
}

func pickClass(STR, INT, WIS, DEX, CON, CHA int) string {
	if CON >= 9 {
		if DEX >= 9 {
			return "Halfling"
		}
		return "Dwarf"
	}
	if INT >= 9 {
		return "Elf"
	}
	classes := []string{"Cleric", "Fighter", "Magic-User", "Thief"}
	return classes[roll(len(classes))-1]
}

func generateTitle(class string, level int) string {
	switch class {
	case "Cleric":
		switch level {
		case 1:
			return "Acolyte"
		case 2:
			return "Adept"
		case 3:
			return "Priest"
		}
	case "Dwarf":
		switch level {
		case 1:
			return "Dwarven Veteran"
		case 2:
			return "Dwarven Warrior"
		case 3:
			return "Dwarven Swordmaster"
		}
	case "Elf":
		switch level {
		case 1:
			return "Veteran-Medium"
		case 2:
			return "Warrior-Seer"
		case 3:
			return "Swordmaster-Conjurer"
		}
	case "Fighter":
		switch level {
		case 1:
			return "Veteran"
		case 2:
			return "Warrior"
		case 3:
			return "Swordmaster"
		}
	case "Halfling":
		switch level {
		case 1:
			return "Halfling Veteran"
		case 2:
			return "Halfling Warrior"
		case 3:
			return "Halfling Swordmaster"
		}
	case "Magic-User":
		switch level {
		case 1:
			return "Medium"
		case 2:
			return "Seer"
		case 3:
			return "Conjurer"
		}
	case "Thief":
		switch level {
		case 1:
			return "Apprentice"
		case 2:
			return "Footpad"
		case 3:
			return "Robber"
		}
	}
	return ""
}

var xpTable = map[string][]int{
	"Cleric":     {1500, 3000, 6000},
	"Dwarf":      {2200, 4400, 8800},
	"Elf":        {4000, 8000, 16000},
	"Fighter":    {2000, 4000, 8000},
	"Halfling":   {2000, 4000, 8000},
	"Magic-User": {2500, 5000, 10000},
	"Thief":      {1200, 2400, 4800},
}

func (c *Character) NextLevel() int {
	if c.Level >= len(xpTable[c.Class]) {
		return -1 // max level reached
	}
	return xpTable[c.Class][c.Level]
}

func (c *Character) PrimeRequisite() int {
	switch c.Class {
	case "Cleric":
		return primeRequisiteSingle(c.WIS)
	case "Dwarf":
		return primeRequisiteSingle(c.STR)
	case "Elf":
		if c.INT >= 13 && c.STR >= 13 {
			if c.INT >= 16 {
				return 10
			}
			return 5
		}
		return 0
	case "Fighter":
		return primeRequisiteSingle(c.STR)
	case "Halfling":
		if c.DEX >= 13 && c.STR >= 13 {
			if c.STR >= 13 {
				return 5
			}
			if c.STR >= 13 && c.DEX >= 13 {
				return 10
			}
		} else if c.STR >= 13 {
			return 5
		}
		return 0
	case "Magic-User":
		return primeRequisiteSingle(c.INT)
	case "Thief":
		return primeRequisiteSingle(c.DEX)
	default:
		return 0
	}
}

func primeRequisiteSingle(prime int) int {
	switch {
	case prime >= 3 && prime <= 5:
		return -20
	case prime >= 6 && prime <= 8:
		return -10
	case prime >= 9 && prime <= 12:
		return 0
	case prime >= 13 && prime <= 15:
		return 5
	case prime >= 16 && prime <= 18:
		return 10
	default:
		return 0
	}
}

var armors = map[int][]string{
	1: {"Leather armor (AC 7)"},
	2: {"Leather armor (AC 7)", "Shield (+1 AC)"},
	3: {"Chainmail (AC 5)"},
	4: {"Chainmail (AC 5)", "Shield (+1 AC)"},
	5: {"Plate armor (AC 3)"},
	6: {"Plate armor (AC 3)", "Shield (+1 AC)"},
}

var weapons = map[int][]string{
	1:  {"Battle axe - 1d8, melee, slow, 2H"},
	2:  {"Crossbow - 1d4, missile (5'-80'/81'-160'/161'-240'), reload, slow, 2H", "20 bolts"},
	3:  {"Hand axe - 1d6, melee, missile (5'-10'/11'-20'/21'-30')"},
	4:  {"Mace - 1d6, blunt, melee"},
	5:  {"Pole arm - 1d10, brace, melee, slow, 2H"},
	6:  {"Short bow - 1d6, missile (5'-50'/51'-100'/101'-150), 2H", "20 arrows"},
	7:  {"Short sword - 1d6, melee"},
	8:  {"Silver dagger - 1d4, melee, missile (5'-10'/11'-20'/21'-30')"},
	9:  {"Sling - 1d4, blunt, missile (5'-40'/41'-80'/81'-160')", "20 stones"},
	10: {"Spear - 1d6, brace, melee, missile (5'-20'/21'-40'/41'-60')"},
	11: {"Sword - 1d8, melee"},
	12: {"War hammer - 1d6, blunt, melee"},
}

var weaponsCleric = map[int][]string{
	1: {"Mace - 1d6, blunt, melee"},
	2: {"Sling - 1d4, blunt, missile (5'-40'/41'-80'/81'-160')", "20 stones"},
	3: {"Staff - 1d4, blunt, melee, 2H"},
	4: {"War hammer - 1d6, blunt, melee"},
}

var gears = map[int][]string{
	1:  {"Crowbar"},
	2:  {"Hammer", "12 iron spikes"},
	3:  {"Holy water"},
	4:  {"Lantern", "3 flasks of oil"},
	5:  {"Mirror (hand-sized, steel)"},
	6:  {"10' Pole"},
	7:  {"50' of Rope"},
	8:  {"50' of Rope", "Grappling hook"},
	9:  {"Large sack"},
	10: {"Small sack"},
	11: {"Stakes (3)", "Mallet"},
	12: {"Wolfsbane (1 bunch)"},
}

func generateInventory(class string) []string {
	armor := armors[roll(6)-1]

	weapon := make([]string, 0, 2)
	switch class {
	case "Cleric":
		weapon = append(weapon, weaponsCleric[roll(4)-1]...)
		weapon = append(weapon, weaponsCleric[roll(4)-1]...)
	case "Magic-User":
		weapon = append(weapon, "Dagger - 1d4, melee, missile (5'-10'/11'-20'/21'-30')")
	default:
		weapon = append(weapon, weapons[roll(12)-1]...)
		weapon = append(weapon, weapons[roll(12)-1]...)
	}

	gear := make([]string, 0, 3)
	gear = append(gear, gears[roll(12)-1]...)
	gear = append(gear, gears[roll(12)-1]...)
	switch class {
	case "Cleric":
		gear = append(gear, "Holy symbol")
	case "Thief":
		gear = append(gear, "Thieves' tools")
	}

	return append(armor, append(weapon, gear...)...)
}
