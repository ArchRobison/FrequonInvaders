package universe

import (
	"github.com/ArchRobison/FrequonInvaders/coloring"
	"github.com/ArchRobison/FrequonInvaders/sound"
	"math/rand"
)

var compressX, compressY float32

var scheme coloring.SchemeBits

func bitCount(s coloring.SchemeBits) int {
	k := 0
	for ; s != 0; k++ {
		s &= s - 1
	}
	return k
}

// At certain scores, certain kinds of damages become possible.
// Here are the thresholds:
//
//		0:  1 stationary alien
//		1:  half-speed alien
//		2:  full-speed alien
//		4:  2 stationary aliens
//		8:  phase, real, or imaginary lost, 1 stationary alien
//		12: 3 aliens
//		16: lose 1 color bit
//		20: 4 aliens
// 		24: lose 2 color bits
//		32: compress x or y
//
// n is score
func setDifficulty(n int) {
	const simpleScheme = coloring.AllBits

	if n == 0 {
		nLiveMax = 1
		velocityMax = 0
		compressX = 1
		compressY = 1
		scheme = simpleScheme
		return
	}

	if n >= 64 {
		gameState = GameWin
		return
	}

	// Set nLiveMax, which is the max number of aliens simultaneously running.
	var liveLimit int
	if n < 4 {
		liveLimit = 1
	} else {
		liveLimit = (n-4)/7 + 2
	}
	if nLiveMax < liveLimit {
		if rand.Intn(2) == 1 {
			nLiveMax++
			// The game is now harder.  Simplify other things that make the game not quite so hard.
			if n < 8 {
				velocityMax = 0
			} else {
				velocityMax /= 2
			}
			if scheme != simpleScheme {
				scheme = simpleScheme
				sound.Play(sound.Bell, 1.5)
			}
			return
		}
	}

	// See if velocityMax should be increased
	// FIXME - should be scaled to screen size
	var velocityLimit float32
	if n <= 4 {
		velocityLimit = float32(n * 15)
	} else {
		velocityLimit = 60
	}
	if velocityMax < velocityLimit {
		if rand.Intn(2) == 1 {
			velocityMax += velocityLimit * 0.5
			if velocityMax > velocityLimit {
				velocityMax = velocityLimit
			}
			if nLiveMax > 1 {
				// Ease up on live limit
				nLiveMax--
			}
			// The game is now harder.
			return
		}
	}

	// Now consider a change in the radar scheme.
	if n >= 8 {
		if rand.Intn(2) == 1 {
			s := scheme
			if s&coloring.CoordinateBits == coloring.CoordinateBits {
				// Break the radar by removing phase, real, or imaginary information.
				switch rand.Intn(3) {
				case 0:
					s &= ^coloring.PhaseBit
				case 1:
					s &= ^coloring.RealBit
				case 2:
					s &= ^coloring.ImagBit
				}
				// Play sound announcing damage.
				sound.Play(sound.Broken, 0.5)
			} else {
				// Fix the radar by restoring all coordinate information (but not color)
				s |= coloring.CoordinateBits
				sound.Play(sound.Bell, 1)
			}
			// Commit the new scheme
			scheme = s
			if s&coloring.CoordinateBits != coloring.CoordinateBits {
				// Game is harder.
				return
			}
		}
	}

	// Consider changing color failures.
	if n >= 16 {
		r := rand.Intn(6)
		if r&1 == 0 {
			var minOnes int
			if n >= 24 {
				minOnes = 1
			} else {
				minOnes = 2
			}
			s := scheme ^ coloring.RedBit<<uint(r%3)
			if bitCount(s&coloring.ColorBits) >= minOnes {
				scheme = s
				// Game is harder
				sound.Play(sound.Broken, 1)
				return
			}
		}
	}

	// Consider compression
	if n >= 32 {
		compressed := false
		switch rand.Intn(2) {
		case 0:
			if compressX > 1./16. {
				compressX *= 0.5
				compressed = true
			}
		case 1:
			if compressY > 1./16. {
				compressY *= 0.5
				compressed = true
			}
		}
		if compressed {
			// Ease up on lost colors
			// FIXME - consider restoring only one color bit
			scheme |= coloring.ColorBits
			if nLiveMax > 1 {
				// Ease up on live limit
				nLiveMax--
			}
			sound.Play(sound.Broken, 0.75)
			return
		}
	}
}

func BoxFraction() (fracX, fracY float32) {
	fracX = compressX
	fracY = compressY
	return
}

func Scheme() coloring.SchemeBits {
	return scheme
}

func SetBoxFraction(frac float32) {
	if frac < 0 || frac > 1 {
		panic("SetBoxFraction: bad frac")
	}
	compressX = frac
	compressY = frac
}

func NKill() int {
	return nKill
}

func BeginGame(isPractice_ bool) {
	SetNLiveMax(0)
	nKill = 0

	if isPractice {
		// Leave visibility to however user currently has it set.
	} else {
		showAlways = false
	}
	isPractice = isPractice_
	gameState = GameActive
}

type GameState int8

var gameState GameState

const (
	GameActive = GameState(iota)
	GameLose
	GameWin
)
