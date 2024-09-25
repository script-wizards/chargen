package character

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/script-wizards/chargen/internal/dice"
)

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
	Gold         int
	Inventory    []string
	Abilities    []string
	ArmorClass   int
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func roll4d6kh3() int {
	rolls := []int{dice.Roll(6), dice.Roll(6), dice.Roll(6), dice.Roll(6)}
	sort.Ints(rolls)
	return rolls[1] + rolls[2] + rolls[3]
}

func rollStats() []int {
	STR := roll4d6kh3()
	INT := roll4d6kh3()
	WIS := roll4d6kh3()
	DEX := roll4d6kh3()
	CON := roll4d6kh3()
	CHA := roll4d6kh3()
	return []int{STR, INT, WIS, DEX, CON, CHA}
}

func unpack(src []int, dst ...*int) {
	for ind, val := range dst {
		*val = src[ind]
	}
}

func isValidClass(class string, stats []int) bool {
	switch strings.ToLower(class) {
	case "dwarf":
		if stats[4] >= 9 {
			return true
		}
	case "elf":
		if stats[1] >= 9 {
			return true
		}
	case "halfling":
		if stats[4] >= 9 && stats[3] >= 9 {
			return true
		}
	default:
		for _, stat := range stats {
			if stat >= 3 && stat <= 18 {
				return true
			}
		}
	}
	return false
}

func rerollUntilValidClass(class string) []int {
	for i := 0; i < 10; i++ {
		stats := rollStats()
		if isValidClass(class, stats) {
			return stats
		}
	}
	return []int{13, 13, 13, 13, 13, 13}
}

func NewCharClass(class string) *Character {
	var STR, INT, WIS, DEX, CON, CHA int
	stats := rerollUntilValidClass(class)
	unpack(stats, &STR, &INT, &WIS, &DEX, &CON, &CHA)

	saves := calcSaves(class)
	hitDie := calcHD(class)

	return &Character{
		Class:        class,
		Level:        1,
		Title:        generateTitle(class, 1),
		Alignment:    alignment(),
		STR:          STR,
		INT:          INT,
		WIS:          WIS,
		DEX:          DEX,
		CON:          CON,
		CHA:          CHA,
		Gold:         dice.Roll3d6(),
		Inventory:    generateInventory(class),
		Abilities:    classAbilities(class),
		SaveDeath:    saves[0],
		SaveWands:    saves[1],
		SaveParalyze: saves[2],
		SaveBreath:   saves[3],
		SaveSpells:   saves[4],
		HitDie:       hitDie,
		HitPoints:    dice.Roll(hitDie),
	}
}

func NewRandomChar() *Character {
	var STR, INT, WIS, DEX, CON, CHA int
	unpack(rollStats(), &STR, &INT, &WIS, &DEX, &CON, &CHA)

	class := pickClass(STR, INT, WIS, DEX, CON, CHA)
	saves := calcSaves(class)
	hitDie := calcHD(class)
	return &Character{
		Class:        class,
		Level:        1,
		Title:        generateTitle(class, 1),
		Alignment:    alignment(),
		STR:          STR,
		INT:          INT,
		WIS:          WIS,
		DEX:          DEX,
		CON:          CON,
		CHA:          CHA,
		Gold:         dice.Roll3d6(),
		Inventory:    generateInventory(class),
		Abilities:    classAbilities(class),
		SaveDeath:    saves[0],
		SaveWands:    saves[1],
		SaveParalyze: saves[2],
		SaveBreath:   saves[3],
		SaveSpells:   saves[4],
		HitDie:       hitDie,
		HitPoints:    dice.Roll(hitDie) + calcMod(CON),
	}
}

func calcHD(class string) int {
	switch class {
	case "cleric", "elf":
		return 6
	case "dwarf", "fighter", "halfling":
		return 8
	case "magic-user", "thief":
		return 4
	default:
		return 0
	}
}

func calcSaves(class string) []int {
	switch class {
	case "cleric":
		return []int{11, 12, 14, 16, 15}
	case "dwarf":
		return []int{8, 9, 10, 13, 12}
	case "elf":
		return []int{12, 13, 13, 15, 15}
	case "fighter":
		return []int{12, 13, 14, 15, 16}
	case "halfling":
		return []int{8, 9, 10, 13, 12}
	case "magic-user":
		return []int{13, 14, 13, 16, 15}
	case "thief":
		return []int{13, 14, 13, 16, 15}
	default:
		return []int{0, 0, 0, 0, 0}
	}
}

func pickClass(STR, INT, WIS, DEX, CON, CHA int) string {
	// Calculate prime requisite bonus for each class
	primeBonuses := map[string]int{
		"cleric":     primeRequisiteSingle(WIS),
		"dwarf":      primeRequisiteSingle(STR),
		"elf":        primeRequisiteElf(INT, STR),
		"fighter":    primeRequisiteSingle(STR),
		"halfling":   primeRequisiteHalfling(DEX, STR),
		"magic-user": primeRequisiteSingle(INT),
		"thief":      primeRequisiteSingle(DEX),
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
	case "cleric":
		return primeRequisiteSingle(WIS)
	case "dwarf":
		return primeRequisiteSingle(STR)
	case "elf":
		return primeRequisiteElf(INT, STR)
	case "fighter":
		return primeRequisiteSingle(STR)
	case "halfling":
		return primeRequisiteHalfling(DEX, STR)
	case "magic-user":
		return primeRequisiteSingle(INT)
	case "thief":
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
	"cleric":     {1500, 3000, 6000},
	"dwarf":      {2200, 4400, 8800},
	"elf":        {4000, 8000, 16000},
	"fighter":    {2000, 4000, 8000},
	"halfling":   {2000, 4000, 8000},
	"magic-user": {2500, 5000, 10000},
	"thief":      {1200, 2400, 4800},
}

func (c *Character) NextLevel() int {
	if c.Level >= len(xpTable[c.Class]) {
		return -1 // max level reached
	}
	return xpTable[c.Class][c.Level-1]
}

func generateTitle(class string, level int) string {
	switch class {
	case "cleric":
		switch level {
		case 1:
			return "Acolyte"
		case 2:
			return "Adept"
		case 3:
			return "Priest"
		}
	case "dwarf":
		switch level {
		case 1:
			return "Dwarven Veteran"
		case 2:
			return "Dwarven Warrior"
		case 3:
			return "Dwarven Swordmaster"
		}
	case "elf":
		switch level {
		case 1:
			return "Veteran-Medium"
		case 2:
			return "Warrior-Seer"
		case 3:
			return "Swordmaster-Conjurer"
		}
	case "fighter":
		switch level {
		case 1:
			return "Veteran"
		case 2:
			return "Warrior"
		case 3:
			return "Swordmaster"
		}
	case "halfling":
		switch level {
		case 1:
			return "Halfling Veteran"
		case 2:
			return "Halfling Warrior"
		case 3:
			return "Halfling Swordmaster"
		}
	case "magic-user":
		switch level {
		case 1:
			return "Medium"
		case 2:
			return "Seer"
		case 3:
			return "Conjurer"
		}
	case "thief":
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

var armorsList = map[int][]string{
	1: {"Leather armor (AC 7)"},
	2: {"Leather armor (AC 7)", "Shield (+1 AC)"},
	3: {"Chainmail (AC 5)"},
	4: {"Chainmail (AC 5)", "Shield (+1 AC)"},
	5: {"Plate armor (AC 3)"},
	6: {"Plate armor (AC 3)", "Shield (+1 AC)"},
}

var weaponsList = map[int][]string{
	1:  {"Battle axe - 1d8, melee, slow, 2H"},
	2:  {"Crossbow - 1d4, missile (80/160/240), reload, slow, 2H", "20 bolts"},
	3:  {"Hand axe - 1d6, melee, missile (10/20/30)"},
	4:  {"Mace - 1d6, blunt, melee"},
	5:  {"Pole arm - 1d10, brace, melee, slow, 2H"},
	6:  {"Short bow - 1d6, missile (50/100'/150), 2H", "20 arrows"},
	7:  {"Short sword - 1d6, melee"},
	8:  {"Silver dagger - 1d4, melee, missile (10/20/30)"},
	9:  {"Sling - 1d4, blunt, missile (40/80/160)", "20 stones"},
	10: {"Spear - 1d6, brace, melee, missile (20/40/60)"},
	11: {"Sword - 1d8, melee"},
	12: {"War hammer - 1d6, blunt, melee"},
}

var weaponsCleric = map[int][]string{
	1: {"Mace - 1d6, blunt, melee"},
	2: {"Sling - 1d4, blunt, missile (40/80/160)", "20 stones"},
	3: {"Staff - 1d4, blunt, melee, 2H"},
	4: {"War hammer - 1d6, blunt, melee"},
}

var gearsList = map[int][]string{
	1:  {"Crowbar"},
	2:  {"Hammer", "12 iron spikes"},
	3:  {"Holy water"},
	4:  {"Lantern", "Oil flask x3 □ □ □"},
	5:  {"Mirror (hand-sized, steel)"},
	6:  {"10' Pole"},
	7:  {"50' of Rope"},
	8:  {"50' of Rope", "Grappling hook"},
	9:  {"Large sack"},
	10: {"Small sack"},
	11: {"Stakes x3 □ □ □", "Mallet"},
	12: {"Wolfsbane (1 bunch)"},
}

func uniqueNumbers(n, max int) []int {
	if n > max {
		return nil
	}

	nums := make([]int, n)
	used := make(map[int]bool)

	for i := 0; i < n; i++ {
		num := rand.Intn(max)
		for used[num] {
			num = rand.Intn(max)
		}
		nums[i] = num
		used[num] = true
	}

	return nums
}

func generateInventory(class string) []string {
	armor := armorsList[dice.Roll(6)-1]
	switch class {
	case "magic-user":
		armor = nil
	case "thief":
		armor = []string{"Leather armor (AC 7)"}
	}

	weapon := make([]string, 0, 10)
	switch class {
	case "cleric":
		n := uniqueNumbers(2, len(weaponsCleric))
		weapon = append(weapon, weaponsCleric[n[0]]...)
		weapon = append(weapon, weaponsCleric[n[1]]...)
	case "magic-user":
		weapon = append(weapon, "Dagger - 1d4, melee, missile (10/20/30)")
	default:
		n := uniqueNumbers(2, len(weaponsList))
		weapon = append(weapon, weaponsList[n[0]]...)
		weapon = append(weapon, weaponsList[n[1]]...)
	}

	gears := make([]string, 0, 10)
	n := uniqueNumbers(2, len(gearsList))
	gears = append(gears, gearsList[n[0]]...)
	gears = append(gears, gearsList[n[1]]...)

	switch class {
	case "cleric":
		gears = append(gears, "Holy symbol")
	case "thief":
		gears = append(gears, "Thieves' tools")
	case "magic-user":
		gears = append(gears, "Spellbook: "+randomSpell(class, 1))
	case "elf":
		gears = append(gears, "Spellbook: "+randomSpell(class, 1))
	}

	return append(armor, append(weapon, gears...)...)
}

func (c *Character) Initiative() int {
	initiative := 0
	if c.Class == "halfling" {
		initiative += 1
	}
	switch c.DEX {
	case 3:
		initiative += -2
	case 4, 5, 6, 7, 8:
		initiative += -1
	case 9, 10, 11, 12:
		initiative += 0
	case 13, 14, 15, 16, 17:
		initiative += 1
	case 18:
		initiative += 2
	default:
		initiative += 0
	}
	return initiative
}

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

var alignments = []string{
	"Lawful",
	"Neutral",
	"Chaotic",
}

func alignment() string {
	return alignments[rng.Intn(3)]
}

func classAbilities(class string) []string {
	return abilitiesMap[class]
}

// TODO: Might want to save line splitting for main.go
var abilitiesMap = map[string][]string{
	"cleric": {
		"Cannot use sharp or piercing weapons, must carry holy symbol",
		"Divine Magic: cast spells (lvl 2+), use divine scrolls/items",
		"Magical Research: Spend time/money to create spells/effects",
		"Turn Undead: 2d6 vs HD (skeleton 7, zombie 9, ghoul 11)",
	},
	"dwarf": {
		"Only use small/normal weapons, can't use longbow or 2H sword",
		"Detect new construction, sliding walls, sloping passages: 2/6",
		"Listen at Doors: 2/6, Infravision: 60'",
		"Languages: Dwarvish, Gnomish, Goblin, Kobold",
	},
	"elf": {
		"Arcane Magic: cast spells, use arcane scrolls/items",
		"Detect hidden/secret doors, listen at doors: 2/6",
		"Infravision: 60, Immune to Ghoul Paralysis",
		"Languages: Elvish, Gnoll, Hobgoblin, Orcish",
	},
	"fighter": {
		"Stronghold: Can build castle or stronghold at any level",
		"",
		"",
		"",
	},
	"halfling": {
		"Only use small weapons/armor, can't use longbow or 2h sword",
		"Hide: woods/undergrowth (90%), still in shadows/cover (2/6)",
		"-2 to AC vs large opponents, +1 initiative, +1 missile attacks",
		"Stronghold: Can build shire at any level",
	},
	"magic-user": {
		"Can only use daggers, can't use shields or wear armor",
		"Arcane Magic: cast spells, use arcane scrolls/items",
		"Magical Research: Spend time/money to create spells/effects",
		"",
	},
	"thief": {
		"Back-Stab: +2 hit/2x dmg attacking unaware enemy from behind",
		"Combat: Cannot wear armor heavier than leather or use shields",
		"CS  TR  HN  HS  MS  OL  PP", // TODO: Thief skills by level
		"87  10   2  10  20  15  20",
	},
}

var spellsArcane = map[int][]string{
	1: {
		"Charm Person",
		"Detect Magic",
		"Floating Disc",
		"Hold Portal",
		"Light",
		"Magic Missile",
		"Protection from Evil",
		"Read Languages",
		"Read Magic",
		"Shield",
		"Sleep",
		"Ventriloquism",
	},
}

func randomSpell(class string, level int) string {
	if class == "elf" || class == "magic-user" {
		spellList := spellsArcane[level]
		return spellList[rng.Intn(len(spellList))]
	}
	return ""
}
