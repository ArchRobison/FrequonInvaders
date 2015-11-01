package universe

import (
	"fmt"
	"github.com/ArchRobison/FrequonInvaders/sound"
	"github.com/ArchRobison/Gophetica/math32"
	"math/rand"
)

type Critter struct {
	Sx, Sy    float32    // Position of particle (units=pixels) - frequency on "fourier" view
	vx, vy    float32    // Velocity of particle (units=pixels/sec)
	Amplitude float32    // Between 0 and 1 - horizontal position on "fall" view and amplitude in "fourier" view.  -1 for "self"
	Progress  float32    // Vertical position on "fall" view, scaled [0,1]
	fallRate  float32    // rate of maturation in maturity/sec
	health    healthType // Initially initialHealth.  Subtracted down to 0. Negative values are death sequence values.  Jumps down to finalHealth at end of sequence
	Show      bool       // If true, show in space domain
	Id        int8       // Index into pastels
}

type healthType int16

const (
	initialHealth healthType = 0x7FFF
	// FIXME - decouple the death animation from the unverse model.
	// Then a negative Amplitude can denote death and -1 can denote a dying Critter
	deathThreshold healthType = -0x8000
)

const killTime = 0.1 // Time it takes being close to kill

const amplitudeDieTime = 2.0 // Time in sec for Frequon to die at full amplitude.

const MaxCritter = 14 // Maximum allowed critters (including self)

var zooStorage [MaxCritter]Critter

var Zoo []Critter

func Init(width, height int32) {
	xSize, ySize = float32(width), float32(height)
	Zoo = zooStorage[:1]
	for k := range zooStorage {
		zooStorage[k].Id = int8(k)
	}
	velocityMax = 0
	// Original sources used (ySize/32) for the square-root of the kill radius.
	// The formula here gives about same answer for widescreen monitors,
	// while adjusting more sensible for other aspect ratios.
	killRadius2 = (xSize * ySize) * ((9. / 16.) / (32 * 32))
}

// Update advances the universe forwards by time interval dt,
// using (selfX,selfY) as the coordinates the player.
func Update(dt float32, selfX, selfY int32) GameState {
	c := Zoo
	if len(c) < 1 {
		panic("universe.Zoo is empty")
	}
	c[0].Sx = float32(selfX)
	c[0].Sy = float32(selfY)
	c[0].Amplitude = -1

	updateLive(dt)
	cullDead()
	tryBirth(dt)
	return gameState
}

var (
	xSize, ySize float32 // Width and height of fourier port, units = pixels
	killRadius2  float32
)

func bounce(sref, vref *float32, limit, dt float32) {
	s, v := *sref, *vref
	s += v * dt
	for {
		if s < 0 {
			s = -s
		} else if s > limit {
			s = 2*limit - s
		} else {
			break
		}
		v = -v
	}
	*sref, *vref = s, v
}

// Update state of live aliens
func updateLive(dt float32) {
	x0, y0 := Zoo[0].Sx, Zoo[0].Sy
	for k := 1; k < len(Zoo); k++ {
		c := &Zoo[k]
		// Update S and v
		bounce(&c.Sx, &c.vx, xSize-1, dt)
		bounce(&c.Sy, &c.vy, ySize-1, dt)
		// Update Progress
		c.Progress += c.fallRate * dt
		// Update health
		if c.health > 0 {
			// Healthy alien
			if c.Progress >= 1 {
				// Alien reached full power!
				if !isPractice {
					gameState = GameLose
				}
				// Mark alien for culling
				c.health = deathThreshold
				continue
			}
			dx := c.Sx - x0
			dy := c.Sy - y0
			if dx*dx+dy*dy <= killRadius2 {
				const killTime = 0.1
				c.health -= healthType(dt * (float32(initialHealth) / killTime))
				if c.health <= 0 {
					c.health = -1 // Transition to death sequence
					sound.Play(sound.Twang, alienPitch[c.Id])
				}
				c.Show = true
			} else {
				c.Show = showAlways
			}
		} else {
			// Dying alien
			c.health -= 1
			c.Show = true
		}
		// Update amplitude
		if c.health > 0 {
			c.Amplitude = math32.Sqrt(c.Progress)
		} else {
			c.Amplitude -= dt * (1 / amplitudeDieTime)
			if c.Amplitude < 0 {
				// Alien croaked.  Mark as dead
				c.health = deathThreshold
				if !isPractice {
					TallyKill()
				}
			}
		}
	}
}

func cullDead() {
	for j := 1; j < len(Zoo); j++ {
		if Zoo[j].health > deathThreshold {
			// Surivor
		} else {
			// Cull corpse by moving Id to end and shrinking slice
			k := len(Zoo) - 1
			deadId := Zoo[j].Id
			Zoo[j] = Zoo[k]
			Zoo[k].Id = deadId
			Zoo = zooStorage[:k]
		}
	}
}

var (
	birthRate   float32 = 1 // units = per second, only an average
	nLiveMax            = 1
	velocityMax float32 = 60 // units = pixels per second
	nKill       int     = 0
	isPractice          = false
)

const π = math32.Pi

func (c *Critter) initAlienVelocity() {
	// Choose random velocity
	// This peculiar 2D distibution was in the original sources.
	// It was probably an error, but is now traditional in Frequon Invaders.
	// It biases towards diagonal directions
	c.vx = math32.Cos(2*π*rand.Float32()) * velocityMax
	c.vy = math32.Sin(2*π*rand.Float32()) * velocityMax
}

// Initialize alien Critter to its birth state.
func (c *Critter) initAlien() {
	// Choose random position
	c.Sx = rand.Float32() * xSize
	c.Sy = rand.Float32() * ySize

	c.initAlienVelocity()

	c.Amplitude = 0

	// Choose random fall rate.  The 0.256 is derived from the original sources.
	c.Progress = 0
	c.fallRate = (rand.Float32() + 1) * 0.0256

	c.health = initialHealth
	c.Show = showAlways
}

var showAlways = false

func SetShowAlways(value bool) {
	showAlways = value
	for k := 1; k < len(Zoo); k++ {
		Zoo[k].Show = showAlways
	}
}

func tryBirth(dt float32) {
	j := len(Zoo)
	if nLiveMax > len(zooStorage)-1 {
		panic(fmt.Sprintf("birth: nLiveMax=%v > len(zooStorage)-1=%v\n", nLiveMax, len(zooStorage)-1))
	}
	if nLive := j - 1; nLive >= nLiveMax {
		// Already reached population limit
		return
	}
	if rand.Float32() > dt*birthRate {
		return
	}

	// Swap in random id
	avail := len(zooStorage) - len(Zoo)
	k := rand.Intn(avail) + j
	zooStorage[j].Id, zooStorage[k].Id = zooStorage[k].Id, zooStorage[j].Id
	Zoo = zooStorage[:j+1]

	// Initialize the alien
	Zoo[j].initAlien()
	sound.Play(sound.AntiTwang, alienPitch[Zoo[j].Id])
}

// FIXME - try to avoid coupling unverse Model to View this way.
func (c *Critter) ImageIndex() int {
	if c.health >= 0 {
		return 0
	} else {
		return -int(c.health)
	}
}

// Embeds five notes of Sprach Zarathustra, but since birth ids are randomly,
// chosen no one will notice
var alienPitch = [MaxCritter]float32{
	3.77549, 1.0000, 1.4983, 2.00000, 2.5198, 2.3784,
	1.1224, 2.9966, 3.36358, 4.0000, 1.2599, 1.3348,
	1.4983, 0.7492}

func init() {
	for _, p := range alienPitch {
		if p < .5 || p > 4 {
			panic("alienPitch wrong")
		}
	}
}

// Set maximum absolute velocity to v pixels per second
func SetVelocityMax(v float32) {
	velocityMax = v
	for k := 1; k < len(Zoo); k++ {
		c := &Zoo[k]
		u := math32.Hypot(c.vx, c.vy)
		if u > 0.001 {
			// Scale velocity to new max
			f := velocityMax / u
			c.vx *= f
			c.vy *= f
		} else {
			// Critter is essentially stationary.  Give it a new velocity.
			c.initAlienVelocity()
		}
	}
}

// Set maximum number of live aliens
func SetNLiveMax(n int) {
	if n < 0 || n > MaxCritter-1 {
		panic(fmt.Sprintf("SetNLiveMax: n=%v", n))
	}
	if n < len(Zoo)-1 {
		// Clip zoo
		Zoo = zooStorage[:n+1]
	}
	nLiveMax = n
}
