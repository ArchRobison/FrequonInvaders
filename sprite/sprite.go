package sprite

import (
	math32 "github.com/ArchRobison/Gophetica/math32"
	nimble "github.com/ArchRobison/Gophetica/nimble"
	"math/rand"
	"sort"
)

type spriteRow struct {
	y int8
	x []int8
}

type fragment struct {
	sx, sy float32 // Position (units = pixels)
	vx, vy float32 // Velocity (units pixels/sec)
}

func makeFragments(radius int, self bool) (frags []fragment) {
	r := float32(radius)
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			d := math32.Hypot(x, y)
			var includePoint bool
			if self {
				// Self is hollow ring
				includePoint = r*0.9 <= d && d <= r*1.1
			} else {
				// Aliens are solid circles
				includePoint = d <= r
			}
			if includePoint {
				// Impart radial velocity to fragment
				var vx, vy, extra float32
				if d > 0 {
					vx, vy = x/d, y/d
					extra = 0.4
				} else {
					vx, vy = 0, 0
					extra = 1.4
				}
				// Random non-radial component
				uy, ux := math32.Sincos(rand.Float32() * (2 * math32.Pi))
				frags = append(frags, fragment{
					x,
					y,
					vx + ux*extra,
					vy + uy*extra,
				})
			}
		}
	}
	return frags
}

func shuffleFragments(frags []fragment) {
	n := int32(len(frags))
	for i := n - 1; i > 0; i-- {
		j := rand.Int31n(i + 1)
		frags[i], frags[j] = frags[j], frags[i]
	}
}

type point struct {
	x, y int8
}

type byYX []point

func (p byYX) Len() int {
	return len(p)
}

func (p byYX) Less(i, j int) bool {
	if p[i].y != p[j].y {
		return p[i].y < p[j].y
	} else {
		return p[i].x < p[j].x
	}
}

func (p byYX) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Sprite struct {
	rows []spriteRow
}

func makeSpriteFromFragments(frags []fragment) (s Sprite) {
	if len(frags) == 0 {
		return
	}
	p := make([]point, 0, len(frags))
	for i := range frags {
		x := int32(math32.Round(frags[i].sx))
		y := int32(math32.Round(frags[i].sy))
		if -128 <= x && x < 128 && -128 <= y && y < 128 {
			p = append(p, point{int8(x), int8(y)})
		}
	}
	sort.Sort(byYX(p))
	// Now build s
	var j int
	for i := 0; i < len(p); i = j {
		y := p[i].y
		for j = i + 1; j < len(p) && p[j].y == y; j++ {
		}
		x := make([]int8, j-i)
		for k := range x {
			x[k] = p[i+k].x
		}
		s.rows = append(s.rows, spriteRow{y, x})
	}
	return
}

// Create sequence of frames
func MakeAnimation(radius int, self bool, seqLen int) (s []Sprite) {
	f := makeFragments(radius, self)
	shuffleFragments(f)
	tScale := float32(len(f)) / float32(seqLen)
	for t := 0; t < seqLen; t++ {
		s = append(s, makeSpriteFromFragments(f[0:len(f)-int(float32(t)*tScale)]))
		for i := range f {
			f[i].sx += f[i].vx
			f[i].sy += f[i].vy
		}
	}
	return
}

func Draw(dst nimble.PixMap, x0, y0 int32, src Sprite, color nimble.Pixel) {
	w, h := dst.Size()
	for _, s := range src.rows {
		y := y0 + int32(s.y)
		if uint32(y) < uint32(h) {
			d := dst.Row(y)
			for _, xoffset := range s.x {
				x := x0 + int32(xoffset)
				if uint32(x) < uint32(w) {
					d[x] = color
				}
			}
		}
	}
}
