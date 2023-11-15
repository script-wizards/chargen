package cairn

var (
	backgrounds = []string{
		"aurifex",
		"barber-surgeon",
		"beast handler",
		"bonekeeper",
		"cutpurse",
		"fieldwarden",
		"fletchwind",
		"foundling",
		"fungal forager",
		"greenwise",
		"half witch",
		"hexenbane",
		"jongleur",
		"kettlewright",
		"marchguard",
		"mountebank",
		"outrider",
		"prowler",
		"rill runner",
		"scrivener",
	}

	names = map[string][]string{
		"aurifex":        {"Hestia", "Basil", "Rune", "Prism", "Ember", "Quintess", "Aludel", "Mordant", "Salaman", "Jazia"},
		"barber-surgeon": {"Wilmot", "Patch", "Lancet", "Sawbones", "Theo", "Cutwell", "Humor", "Landsford", "Goodeye", "Johanna"},
		"beast handler":  {"Amara", "Wulf", "Mireille", "Soren", "Freki", "Aster", "Gerrik", "Boreas", "Delphine", "Matheus"},
		"bonekeeper":     {"Rook", "Ebon", "Moro", "Yew", "Pall", "Leth", "Nix", "Barnaby", "Vesper", "Leder"},
		"cutpurse":       {"Sable", "Lyra", "Eamon", "Salina", "Elara", "Freya", "Isolde", "Sparrow", "Ivy", "Silas"},
		"fieldwarden":    {"Seed", "Thresh", "Dibber", "Sow", "Stalk", "Harrow", "Cobb", "Flax", "Briar", "Rye"},
		"fletchwind":     {"Flint", "Feather", "Crier", "Thunder", "Falcon", "Pluck", "Needle", "Warsong", "Hawk", "Cai"},
		"foundling":      {"Faunus", "Snowdrop", "Wisp", "Silverdew", "Brim", "Solstice", "Steeleye", "Sileas", "Gossamer", "Hazel"},
		"fungal forager": {"Unther", "Woozy", "Hilda", "Current", "Leif", "Ratan", "Mourella", "Lal", "Per", "Madrigal"},
		"greenwise":      {"Briar", "Moss", "Fern", "Lichen", "Root", "Willow", "Sage", "Yarrow", "Rowan", "Ash"},
		"half witch":     {"Solena", "Veles", "Bryn", "Sabine", "Razvan", "Rowena", "Galen", "Nyx", "Vex", "Iwan"},
		"hexenbane":      {"Percival", "Felix", "Isolde", "Wolfram", "Aldric", "Eira", "Oswin", "Ivor", "Brunhilda", "Beatrix"},
		"jongleur":       {"Jax", "Selene", "Baladria", "Ada", "Felix", "Saylor", "Tripp", "Lantos", "Echo", "Jubilo"},
		"kettlewright":   {"Fergus", "Eamon", "Bram", "Idris", "Elara", "Darragh", "Seren", "Rónán", "Berek", "Lorenz"},
		"marchguard":     {"Gann", "Light", "Gale", "Frost", "Thorn", "Reed", "Flint", "Brook", "Brie", "Aasim"},
		"mountebank":     {"Ambrose", "Lucius", "Beauregard", "Cornelius", "Aria", "Seren", "Indigo", "Delphine", "Solene", "Noa"},
		"outrider":       {"Drake", "Cyra", "Keir", "Darius", "Valen", "Rorik", "Yara", "Rui", "Talon", "Jory"},
		"prowler":        {"Winda", "Brielle", "Theron", "Chayse", "Nuja", "Dev", "Raven", "Lyra", "Sable"},
		"rill runner":    {"Gale", "Piper", "Brook", "Adair", "Stone", "Dale", "Wren", "Cliff", "Rain", "Robin"},
		"scrivener":      {"Per", "Stilo", "Akshara", "Pisa", "Ji-Yun", "Kalamos", "Hugo", "Shui", "Kalam", "Julius"},
	}

	traits = [][]string{
		{"athletic", "brawny", "flabby", "lanky", "rugged", "scrawny", "short", "statuesque", "stout", "towering"},
		{"birthmarked", "webbed", "scarred", "marked", "rosy", "soft", "tanned", "tattooed", "tight", "weathered"},
		{"bald", "braided", "curly", "filthy", "frizzy", "long", "luxurious", "oily", "wavy", "wispy"},
		{"bony", "broken", "chiseled", "elongated", "pale", "perfect", "rat-like", "sharp", "square", "sunken"},
		{"blunt", "booming", "cryptic", "droning", "formal", "gravelly", "precise", "squeaky", "stuttering", "whispery"},
		{"antique", "bloody", "elegant", "filthy", "foreign", "frayed", "frumpy", "livery", "rancid", "soiled"},
		{"ambitious", "cautious", "courageous", "disciplined", "gregarious", "honorable", "humble", "merciful", "serene", "tolerant"},
		{"aggressive", "bitter", "craven", "deceitful", "greedy", "lazy", "nervous", "rude", "vain", "vengeful"},
	}
)
