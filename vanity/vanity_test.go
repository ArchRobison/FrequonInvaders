package vanity

import (
	"fmt"
	"math/rand"
	"testing"
)

// Requires visual inspection
func TestIO(t *testing.T) {
	v := []Record{
		Record{Score: 64, Name: "Ice Cream"},
		Record{Score: 50, Name: "Honeydew"},
		Record{Score: 40, Name: "Grape Juice"},
		Record{Score: 30, Name: "Fig"},
		Record{Score: 20, Name: "Elderberry"},
		Record{Score: 10, Name: "Date"},
		Record{Score: 8, Name: "Date"},
		Record{Score: 4, Name: "Cherry"},
		Record{Score: 2, Name: "Banana"},
		Record{Score: 1, Name: "Apple"}}
	for i := 0; i < len(v); i++ {
		WriteToFile(v[:i])
		w, err := ReadFromFile()
		ok := len(w) == i && err == nil
		for j := 0; ok && j < i; j++ {
			ok = v[j] == w[j]
		}
		if !ok {
			fmt.Printf("FAIL w=%v v=%v\n", w, v[0:i])
			t.Fail()
		}
	}
}

func TestInsert(t *testing.T) {
	// Test with small m checks behavior for duplicate scores.
	// Test with large m checks that scores are handled as uint8, not int8.
	for _, m := range []int{3, 256} {
		v := []Record{}
		for i := 0; i < 26; {
			score := uint8(rand.Intn(m))
			if IsWorthyScore(v, score) {
				name := fmt.Sprintf("user %c", 'A'+i)
				v = Insert(v, score, name)
				for j := 0; j < len(v)-1; j++ {
					if v[j].Score < v[j+1].Score || v[j].Score == v[j+1].Score && v[j].Name > v[j+1].Name {
						t.Fail()
					}
				}
				found := false
				for _, r := range v {
					if r.Score == score && r.Name == name {
						found = true
					}
				}
				if !found {
					t.Fail()
				}
				i++
			} else if score != 0 {
				i++
			}
		}
	}
}
