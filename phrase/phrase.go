package phrase

import (
	"fmt"
	"math/rand"
)

// Choose a production for expanding non-terminal variable v.
func (g Grammar) choose(v byte) string {
	rules := g[v]
	return rules[rand.Intn(len(rules))]
}

func (g Grammar) expand(prefix []byte, phrase string, value int) (result []byte) {
	result = prefix
	// Following loop relies on phrase being ASCII
	for i := 0; i < len(phrase); i++ {
		switch phrase[i] {
		case '%':
			i++
			result = g.expand(result, g.choose(phrase[i]), value)
		case '$':
			result = append(result, []byte(fmt.Sprintf("%d", value))...)
		default:
			result = append(result, phrase[i])
		}
	}
	return
}

// phraseBook is a map of grammar productions.
// [k] has productions for non-terminal symbol k
type Grammar map[byte][]string

// Make grammar from array of rules written as strings.
// Each string should have the rule "v: rhs", where v is a singel character
// and rhs is a sequence of characters representating the expansion of v.
// A %u in the rhs is treated as a variable u subject to further expansion.
func MakeGrammar(phraseList []string) (g Grammar) {
	g = make(map[byte][]string)
	for _, p := range phraseList {
		if p[1:3] != ": " {
			panic(fmt.Sprintf("bad grammar rule: %s\n", p))
		}
		key := p[0]
		g[key] = append(g[key], p[3:])
	}
	return
}

// Generate a random phrase starting with given non-terminal root.
func (g Grammar) Generate(root rune) string {
	return g.GenerateWithNumber(root, 0)
}

// Generate a random phrase starting with given non-terminal root,
// and filling number slot with given value.
func (g Grammar) GenerateWithNumber(root rune, value int) string {
	return string(g.expand([]byte{}, g.choose(byte(root)), value))
}
