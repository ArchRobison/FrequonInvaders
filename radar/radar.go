package radar

import (
	"fmt"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
	"math"
)

var (
	frameStorage []nimble.Pixel
	frameSize    int32
	frameHeight  int32
	frameValid   []bool // [k] is true iff frame is valid
)

const nFrame = 120

var (
	frameCounter int32
)

type polarCoor struct {
	d float32 // Distance from origin (not called "r" since that's used for red component)
	φ float32
}

type rgb struct {
	r, g, b float32
}

const π = math.Pi

var (
	xSize, ySize     int32   // width and height of view (units = pixels)
	xScale, yScale   float32 // Scale for mapping view onto [-1,1]x[-1,1]
	xOffset, yOffset float32 // Offsets for mapping view onto [-1,1]x[-1,1]
)

// Linearized arrays of dimension height x width
var (
	polar []polarCoor // Storage for polar coordinates
	clut  []rgb       // Storate for clut
)

// Init specifies the size of the PixMap seen by subsequent calls to Draw
func Init(width, height int32) {
	if width < 0 || height < 0 {
		panic(fmt.Sprintf("radar.Init: width=%v height=%v\n", width, height))
	}
	xSize, ySize = width, height

	xScale = 2 / float32(width-1)
	yScale = 2 / float32(height-1)
	xOffset = 1.0 - xScale*float32(width)
	yOffset = 1.0 - yScale*float32(height)

	// Allocate clut
	clut = make([]rgb, height*width)

	// Allocate frames
	frameSize = height * width
	frameStorage = make([]nimble.Pixel, nFrame*frameSize)
	frameValid = make([]bool, nFrame)

	// Compute polar coordinates
	polar = make([]polarCoor, height*width)
	for i := int32(0); i < height; i++ {
		for j := int32(0); j < width; j++ {
			x := float32(j)*xScale + xOffset
			y := float32(i)*yScale + yOffset
			polar[i*width+j] = polarCoor{d: math32.Hypot(x, y), φ: math32.Atan2(-y, x)}
		}
	}
}

// Get pixels for the kth frame in the animation of the radar.
// Frames are constructed lazily to avoid introducing a big delay when the color scheme changes.
func getFrame(k int32) []nimble.Pixel {
	f := frameStorage[k*frameSize : (k+1)*frameSize]
	if !frameValid[k] {
		// Need to compute the frame
		width, height := xSize, ySize
		// Construct the frame, incorporating "sweep"
		θ := float32(k)/nFrame*(2*π) - π
		for i := int32(0); i < height; i++ {
			for j := int32(0); j < width; j++ {
				var color nimble.Pixel
				p := polar[i*width+j]
				if p.d <= 1 {
					delta := p.φ - θ
					if delta < 0 {
						delta += 2 * π
					}
					factor := delta * (1 / (2 * π))
					c := clut[i*width+j]
					color = nimble.RGB(c.r*factor, c.g*factor, c.b*factor)
				} else {
					color = nimble.Black
				}
				f[i*width+j] = color
			}
		}
		frameValid[k] = true
	}
	return f
}

type colorMap interface {
	Color(x, y float32) (r, g, b float32)
}

var currentMap colorMap

func setColoring(cm colorMap) {
	if cm == currentMap {
		return
	}
	currentMap = cm

	// Compute clut as if there was no "radar sweep"
	width, height := xSize, ySize
	for i := int32(0); i < height; i++ {
		for j := int32(0); j < width; j++ {
			x := float32(j)*xScale + xOffset
			y := float32(i)*yScale + yOffset
			r, g, b := cm.Color(x, y)
			clut[i*width+j] = rgb{r, g, b}
		}
	}

	// Mark frames a invalid
	for i := range frameValid {
		frameValid[i] = false
	}
}

func Draw(pm nimble.PixMap, cm colorMap, running bool) {
	setColoring(cm)
	width, height := pm.Size()
	if width != xSize || height != ySize {
		panic(fmt.Sprintf("radar.Draw: (width,height)=(%v,%v) (xSize,ySize)=(%v,%v)\n",
			width, height, xSize, ySize))
	}
	src := nimble.MakePixMap(width, height, getFrame(frameCounter), width)
	pm.Copy(0, 0, &src)
	if running {
		frameCounter = (frameCounter + 1) % (int32(len(frameStorage)) / frameSize)
	}
}
