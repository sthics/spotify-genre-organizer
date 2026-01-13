package genres

import "strings"

var ParentGenres = []string{
	"Rock", "Pop", "Hip-Hop", "Electronic", "R&B", "Jazz", "Classical",
	"Country", "Metal", "Folk", "Latin", "Blues", "Reggae", "Punk",
	"Indie", "Soul", "Funk", "World", "Other",
}

// GenrePriority defines tie-breaking order (earlier = higher priority)
// More specific genres come before broader ones
var GenrePriority = []string{
	"Classical", "Jazz", "Blues", "Reggae", "Folk", "Country",
	"Metal", "Punk", "Funk", "Soul", "R&B", "Latin", "World",
	"Rock", "Electronic", "Hip-Hop", "Pop", "Indie", "Other",
}

var genreMapping = map[string]string{
	// Rock
	"rock": "Rock", "indie rock": "Rock", "alternative rock": "Rock",
	"garage rock": "Rock", "classic rock": "Rock", "hard rock": "Rock",
	"soft rock": "Rock", "progressive rock": "Rock", "psychedelic rock": "Rock",
	"art rock": "Rock", "glam rock": "Rock", "grunge": "Rock",
	"post-rock": "Rock", "shoegaze": "Rock", "britpop": "Rock",

	// Pop
	"pop": "Pop", "indie pop": "Pop", "synth-pop": "Pop", "electropop": "Pop",
	"dance pop": "Pop", "art pop": "Pop", "dream pop": "Pop",
	"chamber pop": "Pop", "power pop": "Pop", "teen pop": "Pop",
	"k-pop": "Pop", "j-pop": "Pop",

	// Hip-Hop
	"hip hop": "Hip-Hop", "rap": "Hip-Hop", "trap": "Hip-Hop",
	"conscious hip hop": "Hip-Hop", "gangsta rap": "Hip-Hop",
	"underground hip hop": "Hip-Hop", "boom bap": "Hip-Hop",
	"drill": "Hip-Hop", "crunk": "Hip-Hop", "grime": "Hip-Hop",

	// Electronic
	"electronic": "Electronic", "edm": "Electronic", "house": "Electronic",
	"techno": "Electronic", "trance": "Electronic", "dubstep": "Electronic",
	"drum and bass": "Electronic", "ambient": "Electronic", "idm": "Electronic",
	"downtempo": "Electronic", "trip hop": "Electronic", "chillwave": "Electronic",
	"synthwave": "Electronic", "deep house": "Electronic", "tech house": "Electronic",
	"progressive house": "Electronic",

	// R&B
	"r&b": "R&B", "rnb": "R&B", "contemporary r&b": "R&B",
	"neo soul": "R&B", "new jack swing": "R&B", "quiet storm": "R&B",

	// Jazz
	"jazz": "Jazz", "jazz fusion": "Jazz", "smooth jazz": "Jazz",
	"bebop": "Jazz", "cool jazz": "Jazz", "free jazz": "Jazz",
	"acid jazz": "Jazz", "nu jazz": "Jazz", "swing": "Jazz", "big band": "Jazz",

	// Classical
	"classical": "Classical", "baroque": "Classical", "romantic": "Classical",
	"contemporary classical": "Classical", "opera": "Classical",
	"orchestral": "Classical", "chamber music": "Classical", "symphony": "Classical",

	// Country
	"country": "Country", "country rock": "Country", "alt-country": "Country",
	"bluegrass": "Country", "americana": "Country", "outlaw country": "Country",
	"country pop": "Country",

	// Metal
	"metal": "Metal", "heavy metal": "Metal", "thrash metal": "Metal",
	"death metal": "Metal", "black metal": "Metal", "doom metal": "Metal",
	"power metal": "Metal", "progressive metal": "Metal", "nu metal": "Metal",
	"metalcore": "Metal",

	// Folk
	"folk": "Folk", "indie folk": "Folk", "folk rock": "Folk",
	"freak folk": "Folk", "contemporary folk": "Folk", "traditional folk": "Folk",

	// Latin
	"latin": "Latin", "reggaeton": "Latin", "salsa": "Latin",
	"bachata": "Latin", "cumbia": "Latin", "bossa nova": "Latin",
	"latin pop": "Latin", "latin rock": "Latin",

	// Blues
	"blues": "Blues", "electric blues": "Blues", "delta blues": "Blues",
	"chicago blues": "Blues", "blues rock": "Blues",

	// Reggae
	"reggae": "Reggae", "dub": "Reggae", "ska": "Reggae",
	"dancehall": "Reggae", "roots reggae": "Reggae",

	// Punk
	"punk": "Punk", "punk rock": "Punk", "pop punk": "Punk",
	"post-punk": "Punk", "hardcore punk": "Punk", "emo": "Punk", "skate punk": "Punk",

	// Indie
	"indie": "Indie", "lo-fi": "Indie", "bedroom pop": "Indie",

	// Soul
	"soul": "Soul", "motown": "Soul", "northern soul": "Soul", "southern soul": "Soul",

	// Funk
	"funk": "Funk", "p-funk": "Funk", "funk rock": "Funk", "disco": "Funk",

	// World
	"world": "World", "afrobeat": "World", "afropop": "World",
	"celtic": "World", "flamenco": "World", "indian": "World", "middle eastern": "World",
}

func ConsolidateGenre(microGenre string) string {
	normalized := strings.ToLower(strings.TrimSpace(microGenre))

	if parent, ok := genreMapping[normalized]; ok {
		return parent
	}

	for micro, parent := range genreMapping {
		if strings.Contains(normalized, micro) || strings.Contains(micro, normalized) {
			return parent
		}
	}

	return "Other"
}

func GetParentGenres() []string {
	return ParentGenres
}

func ConsolidateGenres(microGenres []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, g := range microGenres {
		parent := ConsolidateGenre(g)
		if !seen[parent] {
			seen[parent] = true
			result = append(result, parent)
		}
	}

	return result
}

// ScoreGenres takes all artist genres and returns the best-fit parent genre
// using weighted voting with priority-based tie-breaking
func ScoreGenres(microGenres []string) string {
	if len(microGenres) == 0 {
		return "Other"
	}

	// Count votes for each parent genre
	votes := make(map[string]int)
	for _, g := range microGenres {
		parent := ConsolidateGenre(g)
		votes[parent]++
	}

	// Find max vote count
	maxVotes := 0
	for _, count := range votes {
		if count > maxVotes {
			maxVotes = count
		}
	}

	// Collect all genres with max votes
	var candidates []string
	for genre, count := range votes {
		if count == maxVotes {
			candidates = append(candidates, genre)
		}
	}

	// If single winner, return it
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Break tie using priority order
	for _, priorityGenre := range GenrePriority {
		for _, candidate := range candidates {
			if candidate == priorityGenre {
				return candidate
			}
		}
	}

	return "Other"
}
