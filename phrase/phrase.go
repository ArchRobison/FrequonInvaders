package phrase

import (
	"math/rand"
)

var phraseList = []string{
	"1: Boot %A %O",
	"1: Reset %A %O",
	"1: Load %A %O",
	"2: Enable %A %O",
	"2: Engage %A %O",
	"2: Digitize %A %O",
	"3: Toggle %A %O",
	"3: Activate %A %O",
	"4: Invert %A %O",
	"4: Turn on %A %O",
	"4: Spin up %A %O",
	"5: Tune %A %O",
	"5: Rev up %A %O",
	"6: Switch on %A %O",
	"6: Check %A %O",
	"6: Lock %A %O",
	"7: Power on %A %O",
	"7: Energize %A %O",
	"8: Compress %A %O",
	"O: overthruster",
	"O: inverter",
	"O: integrator",
	"O: differentiator",
	"O: accumulator",
	"O: accelerator",
	"O: processors",
	"O: twistor",
	"O: eigenvector",
	"O: function",
	"O: reactor",
	"O: transporter",
	"A: micro",
	"A: elliptic",
	"A: interprocedural",
	"A: hyperfinite",
	"A: photonic",
	"A: quantum",
	"A: nano",
	"A: alpha",
	"A: beta",
	"A: psi",
	"A: transwarp",
	"A: harmonic",
	"A: complex",
	"A: hyperbolic",
	"A: pseudomorphic",
	"A: acoustic",
	"S: Score = @S",
	"H: Your score of @S merits recording.",
	"W: You Win!",
}

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

func expand(prefix []byte, phrase string) (result []byte) {
	result = prefix
	// Following loop relies on phrase being ASCII
	for i := 0; i < len(phrase); i++ {
		switch phrase[i] {
		case '%':
			result = expand(result, choose(phrase[i+1]))
			i++
		case '@':
			// FIXME
			i++
		default:
			// FIXME - suppress double spaces here
			result = append(result, phrase[i])
		}
	}
	return
}

/*
static char* GeneratePhraseLoop( char type, char * dst ) {
	unsigned count = 0;
	// Broken MSVC++ 6.0 compiler does not follow ISO for scope rules
	Phrase * p;
	for( p = PhraseBook[type]; p; p=p->next ) {
		++count;
	}
	if( count==0 ) {
		*dst++ = '?';
		return dst;
	}
	p = PhraseBook[type];
	for( unsigned k = (NimbleRandom()>>8) % count; k>0; --k ) {
		p = p->next;
	}
	for( const char * src = p->string; int c=*src; ++src ) {
		if( dst>=&Buffer[PHRASE_SIZE_MAX+1] ) return NULL;
		if( c=='%' ) {
			if( !*++src ) break;
			dst = GeneratePhraseLoop( *src, dst );	// Recurse
			if( dst==NULL ) return NULL;
		} else if( c=='@' ) {
			*dst = 0;
			switch( c=*++src ) {
				case 'N':
					// Replace with newline
					*dst = '\n';
					break;
				case 'S': {
					// Replace with score.
					// sprintf is not in basic Mac library, so do this the hard way.
					dst = ToDecimal(dst,TheUniverse.n_kill);
					break;
				}
			}
		} else {
			if( c==' ' && dst[-1]==' ' ) continue;
			*dst++ = c;
		}
	}
	return dst;
}
*/

func Generate(root rune) string {
	return string(expand([]byte{}, choose(byte(root))))
}
