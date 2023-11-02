package character

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

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
	SaveDeath    int
	SaveWands    int
	SaveParalyze int
	SaveBreath   int
	SaveSpells   int
	HitDie       int
	HitPoints    int
	Inventory    []string
	ArmorClass   int
}

func roll(die int) int {
	return 1 + rng.Intn(die)
}

func roll4d6kh3() int {
	rolls := []int{roll(6), roll(6), roll(6), roll(6)}
	sort.Ints(rolls)
	return rolls[1] + rolls[2] + rolls[3]
}

func NewRandomChar() *Character {
	// roll stats first before picking class
	STR := roll4d6kh3()
	INT := roll4d6kh3()
	WIS := roll4d6kh3()
	DEX := roll4d6kh3()
	CON := roll4d6kh3()
	CHA := roll4d6kh3()

	class := pickClass(STR, INT, WIS, DEX, CON, CHA)
	saves := calcSaves(class)
	hitDie := calcHD(class)
	return &Character{
		Class:        class,
		Level:        1,
		Title:        generateTitle(class, 1),
		STR:          STR,
		INT:          INT,
		WIS:          WIS,
		DEX:          DEX,
		CON:          CON,
		CHA:          CHA,
		Inventory:    generateInventory(class),
		SaveDeath:    saves[0],
		SaveWands:    saves[1],
		SaveParalyze: saves[2],
		SaveBreath:   saves[3],
		SaveSpells:   saves[4],
		HitDie:       hitDie,
		HitPoints:    roll(hitDie),
	}
}

func calcHD(class string) int {
	switch class {
	case "Cleric", "Elf":
		return 6
	case "Dwarf", "Fighter", "Halfling":
		return 8
	case "Magic-User", "Thief":
		return 4
	default:
		return 0
	}
}

func calcSaves(class string) []int {
	switch class {
	case "Cleric":
		return []int{11, 12, 14, 16, 15}
	case "Dwarf":
		return []int{8, 9, 10, 13, 12}
	case "Elf":
		return []int{12, 13, 13, 15, 15}
	case "Fighter":
		return []int{12, 13, 14, 15, 16}
	case "Halfling":
		return []int{8, 9, 10, 13, 12}
	case "Magic-User":
		return []int{13, 14, 13, 16, 15}
	case "Thief":
		return []int{13, 14, 13, 16, 15}
	default:
		return []int{0, 0, 0, 0, 0}
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
	// Calculate prime requisite bonus for each class
	primeBonuses := map[string]int{
		"Cleric":     primeRequisiteSingle(WIS),
		"Dwarf":      primeRequisiteSingle(STR),
		"Elf":        primeRequisiteElf(INT, STR),
		"Fighter":    primeRequisiteSingle(STR),
		"Halfling":   primeRequisiteHalfling(DEX, STR),
		"Magic-User": primeRequisiteSingle(INT),
		"Thief":      primeRequisiteSingle(DEX),
	}

	// Look for classes that give the highest bonus
	highestBonus := -100
	var classesWithHighestBonus []string
	for class, bonus := range primeBonuses {
		if bonus > highestBonus {
			highestBonus = bonus
			classesWithHighestBonus = []string{class}
		} else if bonus == highestBonus {
			classesWithHighestBonus = append(classesWithHighestBonus, class)
		}
	}

	// Pick a random class from the ones that give the highest bonus
	if len(classesWithHighestBonus) > 0 {
		return classesWithHighestBonus[rng.Intn(len(classesWithHighestBonus))]
	}

	// Look for classes that give the second highest bonus
	secondHighestBonus := -100
	var classesWithSecondHighestBonus []string
	for class, bonus := range primeBonuses {
		if bonus > secondHighestBonus && bonus < highestBonus {
			secondHighestBonus = bonus
			classesWithSecondHighestBonus = []string{class}
		} else if bonus == secondHighestBonus {
			classesWithSecondHighestBonus = append(classesWithSecondHighestBonus, class)
		}
	}

	// Pick a random class from the ones that give the second highest bonus
	if len(classesWithSecondHighestBonus) > 0 {
		return classesWithSecondHighestBonus[rng.Intn(len(classesWithSecondHighestBonus))]
	}

	// Look for classes that give 0 bonus
	var classesWithZeroBonus []string
	for class, bonus := range primeBonuses {
		if bonus == 0 {
			classesWithZeroBonus = append(classesWithZeroBonus, class)
		}
	}

	// Pick a random class from the ones that give 0 bonus
	if len(classesWithZeroBonus) > 0 {
		return classesWithZeroBonus[rng.Intn(len(classesWithZeroBonus))]
	}

	// Look for classes that give the lowest bonus
	lowestBonus := 100
	var classesWithLowestBonus []string
	for class, bonus := range primeBonuses {
		if bonus < lowestBonus {
			lowestBonus = bonus
			classesWithLowestBonus = []string{class}
		} else if bonus == lowestBonus {
			classesWithLowestBonus = append(classesWithLowestBonus, class)
		}
	}

	// Pick a random class from the ones that give the lowest bonus
	if len(classesWithLowestBonus) > 0 {
		return classesWithLowestBonus[rng.Intn(len(classesWithLowestBonus))]
	}

	// If no class was found, return an empty string
	return ""
}

func PrimeRequisite(class string, STR, INT, WIS, DEX, CON, CHA int) int {
	switch class {
	case "Cleric":
		return primeRequisiteSingle(WIS)
	case "Dwarf":
		return primeRequisiteSingle(STR)
	case "Elf":
		return primeRequisiteElf(INT, STR)
	case "Fighter":
		return primeRequisiteSingle(STR)
	case "Halfling":
		return primeRequisiteHalfling(DEX, STR)
	case "Magic-User":
		return primeRequisiteSingle(INT)
	case "Thief":
		return primeRequisiteSingle(DEX)
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

func primeRequisiteElf(INT, STR int) int {
	if INT >= 13 && STR >= 13 {
		if INT >= 16 {
			return 10
		}
		return 5
	}
	return 0
}

func primeRequisiteHalfling(prime1, prime2 int) int {
	if prime1 >= 13 || prime2 >= 13 {
		if prime1 >= 13 && prime2 >= 13 {
			return 10
		}
		return 5
	}
	return 0
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
	return xpTable[c.Class][c.Level-1]
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
	switch class {
	case "Magic-User":
		armor = nil
	case "Thief":
		armor = []string{"Leather armor (AC 7)"}
	}

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

/*
by default, a character's AC is 9. Go through the inventory and see if it has armor with (AC #) in it. if so, that becomes the new AC. go through the inventory and see if there is a shield. If so, subtract 1 from the AC. finally, subtract the DEX modifier the AC.
*/
func (c *Character) SetAC() int {
	ac := 9
	for _, item := range c.Inventory {
		if strings.Contains(item, "AC ") {
			acStr := strings.Split(item, "AC ")[1][0:1]
			acInt, err := strconv.Atoi(acStr)
			if err == nil {
				ac = acInt
			}
		}
		if strings.Contains(item, "Shield") {
			ac -= 1
		}
	}
	ac -= calcMod(c.DEX)
	c.ArmorClass = ac
	return ac
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
