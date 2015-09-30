package phrase

import (
	"fmt"
	"math/rand"
)

// Map of grammar productions.  [k] has productions for non-terminal symbol k
var phraseBook = make(map[byte][]string)

func choose(root byte) string {
	plist := phraseBook[root]
	return plist[rand.Intn(len(plist))]
}

func init() {
	for _, p := range phraseList {
		key := p[0]
		phraseBook[key] = append(phraseBook[key], p[3:])
	}
}

func expand(prefix []byte, phrase string, value int) (result []byte) {
	result = prefix
	// Following loop relies on phrase being ASCII
	for i := 0; i < len(phrase); i++ {
		switch phrase[i] {
		case '%':
			result = expand(result, choose(phrase[i+1]), value)
			i++
		case '$':
			result = append(result, []byte(fmt.Sprintf("%d", value))...)
		default:
			// FIXME - need to suppress double spaces here?
			result = append(result, phrase[i])
		}
	}
	return
}

func Generate(root rune) string {
	return GenerateWithNumber(root, 0)
}

func GenerateWithNumber(root rune, value int) string {
	return string(expand([]byte{}, choose(byte(root)), value))
}
