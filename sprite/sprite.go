package sprite

import (
	"math/rand"
	. "github.com/ArchRobison/FrequonInvaders/math32"
	. "github.com/ArchRobison/NimbleDraw"
	"sort"
)

type spriteRow struct {
	y int8
	x []int8
}

func drawSprite(dst *PixMap, x0, y0 int32, src []spriteRow, color Pixel) {
	w := uint32(dst.Width())
	h := uint32(dst.Height())
	for _, s := range src {
		y := y0 + int32(s.y)
		if uint32(y) < h {
			d := dst.Row(y)
			for _, xoffset := range s.x {
				x := x0 + int32(xoffset)
				if uint32(x) < w {
					d[x] = color
				}
			}
		}
	}
}

type fragment struct {
	sx, sy float32
	vx, vy float32
}

func makeFragments(radius int, self bool) (frags []fragment) {
	r := float32(radius)
	for y := -r; y <= r; y++ {
		for x := -r; x <= r; x++ {
			d := Hypot(x, y)
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
				θ := rand.Float32() * (2 * Pi)
				frags = append(frags, fragment{
					x,
					y,
					vx + Cos(θ)*extra,
					vy + Sin(θ)*extra,
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

func makeSpriteFromFragments(frags []fragment) (s []spriteRow) {
	if len(frags) == 0 {
		return
	}
	p := make([]point, 0, len(frags))
	for i := range frags {
		x := RoundToInt(frags[i].sx)
		y := RoundToInt(frags[i].sy)
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
		s = append(s, spriteRow{y, x})
	}
	return
}

type Seq [][]spriteRow

func Len(s Seq) int {
	return len(s)
}

func MakeSeq(radius int, self bool, seqLen int) (s [][]spriteRow) {
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

func Draw(pm PixMap, x0, y0 int32, s Seq, t int, color Pixel) {
	drawSprite(&pm, x0, y0, s[t], color)
}
