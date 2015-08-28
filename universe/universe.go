package universe

import (
    "fmt"
    . "github.com/ArchRobison/FrequonInvaders/math32"
    "math/rand"
)

type Critter struct {
	Sx, Sy    float32 // Position of particle (units=pixels) - frequency on "fourier" view
	vx, vy    float32 // Velocity of particle (units=pixels/sec)
	Amplitude float32 // Between 0 and 1 - horizontal position on "fall" view and amplitude in "fourier" view.  -1 for "self"
	Maturity  float32 // Vertical position on "fall" view, scaled [0,1]
    fallRate  float32 // rate of maturation in maturity/sec
    health	  healthType // Initially FullHealth.  Subtracted down to 0. Negative values are death sequence values. 
						// Jumps down to finalHealth at end of sequence  
    id 		  int8	  // Index into pastels
}

type healthType int16

const (
    initialHealth healthType = 0x7FFF
    deathThreshold healthType = -0x8000
)

const sickRate = float32(initialHealth)/10.

const amplitudeDieRate = 2	// units = amplitude per sec

const MaxCritter = 16 // Maximum allowed critters 

var zooStorage [MaxCritter]Critter

var Zoo[]Critter

func Init(width, height int32) {
    xSize, ySize = float32(width), float32(height)
    Zoo = zooStorage[0:1]
    for k:=range zooStorage {
        zooStorage[k].id = int8(k)
    }
    // Original sources used ySize/32 here.  
    // The formula here gives the same answer for widescreen monitors,
    // while adjusting more sensible for other aspect ratios.
    killRadius = Sqrt(xSize*ySize)*(3./128)
}

func Update(dt float32) {
    updateLive(dt)
    cullDead()
    tryBirth(dt)
}

var (
    xSize, ySize float32 // Width and height of fourier port, units = pixels
    killRadius float32
)

func bounce(sref, vref *float32, limit, dt float32) {
	v := *vref
	s := *sref + v*dt
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
	*sref = s
	*vref = v
}

func updateLive(dt float32) {
	x0, y0 := Zoo[0].Sx, Zoo[0].Sy
	for k := range Zoo[1:] {
		c := &Zoo[k]
		// Update S and v
		bounce(&c.Sx, &c.vx, dt, xSize-1)
		bounce(&c.Sy, &c.vy, dt, ySize-1)
		// Update Maturity
		c.Maturity += c.fallRate*dt
		// Update health 
		if c.health > 0 {
			// Healthy alien
			d := Hypot(c.Sx-x0, c.Sy-y0)
			if d < killRadius {
				const killTime = 0.1
				c.health -= healthType(dt * sickRate)
				if c.health <= 0 {
					c.health = -1				// Transition to death sequence
					// FIXME - play sound here
				}
			}
		} else {
            // Dying alien
			c.health -= 1
		}   
		// Update amplitude
		if c.health > 0 {
		    c.Amplitude = Sqrt(c.Maturity)
		} else {
			c.Amplitude -= dt * amplitudeDieRate
		}
	}
}

func cullDead() {
    for j:=0; j<len(Zoo); {
        if Zoo[j].health>deathThreshold {
		    // Surivor
		} else {
		    j++
			// Cull it by moving id to end and shrinking slice
			k := len(Zoo)-1
		    deadId := Zoo[j].id
		    Zoo[j] = Zoo[k]
			Zoo[k].id = deadId
			Zoo = Zoo[:k]
		}
    }
}

var (birthRate float32 = 1	// units = per second, only an average
    maxLive = 1
    velocityMax float32 = 1
)

const π = Pi

func (c * Critter) init() {
    // Choose random position
    c.Sx = rand.Float32()*xSize
    c.Sy = rand.Float32()*ySize

    // Choose random velocity
    // This peculiar 2D distibution was in the original sources.
    // It was probably an error, but now is traditional in Frequon Invaders.
    // It biases towards diagonal directions
    c.vx = Cos( 2*π*rand.Float32() ) * velocityMax
    c.vy = Sin( 2*π*rand.Float32() ) * velocityMax

    c.Amplitude = 0

    // Choose random fall rate.  The 0.256 is derived from the original sources.
    c.Maturity = 0
	c.fallRate = (rand.Float32()+1)*0.0256

    c.health = initialHealth

	// Leave id alone
}

func tryBirth(dt float32) {
    j := len(Zoo)
    if maxLive > cap(Zoo)-1 {
		panic(fmt.Sprintf("birth: maxLIve=%v > cap(Zoo)-1=%v\n", maxLive, cap(Zoo)-1))
    }
    nLive := j-1
    if nLive > maxLive {
		panic(fmt.Sprintf("birth: nLive=%v > maxLive=%v\n", nLive, maxLive))
	}
	if nLive >= maxLive {
		return
	}
	if( rand.Float32()>dt*birthRate ) {
	    return
    }
	// Swap in random id
    avail := cap(Zoo)-len(Zoo)
	k := rand.Intn(avail)+j
	Zoo[j].id, Zoo[k].id = Zoo[k].id, Zoo[j].id

    // Initialize the critter
	Zoo[j].init()
}

