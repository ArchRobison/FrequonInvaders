package radar

import (
	"fmt"
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
	"math"
)

var frameStorage []nimble.Pixel
var frameSize, frameHeight int32

func getFrame(k int32) []nimble.Pixel {
	return frameStorage[k*frameSize : (k+1)*frameSize]
}

const nFrame = 120

var frameCounter int32

type polarCoor struct {
	d float32 // Distance from origin (r not used to avoid confusion with red component)
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

type colorMap interface {
	Color(x, y float32) (r, g, b float32)
}

func SetColoring(cm colorMap) {
	width, height := xSize, ySize

	// Compute clut as if there was no "radar sweep"
	for i := int32(0); i < height; i++ {
		for j := int32(0); j < width; j++ {
			x := float32(j)*xScale + xOffset
			y := float32(i)*yScale + yOffset
			r, g, b := cm.Color(x, y)
			clut[i*width+j] = rgb{r, g, b}
		}
	}
	black := nimble.Gray(0)

	// Construct the frames, incorporating "sweep"
	for t := int32(0); t < nFrame; t++ {
		θ := float32(t)/nFrame*(2*π) - π
		pm := nimble.MakePixMap(width, height, getFrame(t), width)
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
					color = black
				}
				pm.SetPixel(j, i, color)
			}
		}
	}

}

func Draw(pm nimble.PixMap, running bool) {
	width, height := pm.Size()
	if width != xSize || height != ySize {
		panic(fmt.Sprintf("radar.Draw: (width,height)=(%v,%v) (xSize,ySize)=(%v,%v)\n",
			width, height, xSize, ySize))
	}
	src := nimble.MakePixMap(width, height, getFrame(frameCounter), width)
	pm.Copy(0, 0, &src)
	frameCounter = (frameCounter + 1) % (int32(len(frameStorage)) / frameSize)
}
