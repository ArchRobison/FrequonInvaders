package vanity

const (
	recordMax        = 10
	vanityNameLenMax = 20
)

// Record has the score and name for a game.
// Slices of these records are kept in order of descending score.
type Record struct {
	Score uint8
	Name  string
}

// IsWorthyScore returns true if score is worthy of adding to vanity board.
func IsWorthyScore(v []Record, score uint8) bool {
	return score > 0 && (len(v) < recordMax || v[len(v)-1].Score < score)
}

// Insert score/name into vanity
func Insert(v []Record, score uint8, name string) []Record {
	if len(v) >= recordMax {
		v[len(v)-1] = Record{score, name}
	} else {
		v = append(v, Record{score, name})
	}
	// Bubble it up to keep records sorted by score
	for i := len(v) - 1; i > 0; i-- {
		if v[i-1].Score < v[i].Score {
			v[i-1], v[i] = v[i], v[i-1]
		}
	}
	return v
}
