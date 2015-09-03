package sprite

import (
	"fmt"
	"github.com/ArchRobison/Gophetica/nimble"
	"testing"
)

func assert(expr bool, msg string, t *testing.T) {
	if !expr {
		t.Fail()
		panic(msg)
	}
}

// Requires visual inspection of stdout
func TestMakeFragments(t *testing.T) {
	const (
		w = 20
		h = 30
		a = 12
		b = 14
	)
	src := MakeAnimation(5, false, 1)
	dst := nimble.MakePixMap(w, h, make([]nimble.Pixel, w*h), w)
	background := nimble.RGB(0.5, 0.25, 1)
	foreground := nimble.RGB(0.25, 1, 0.25)
	dst.Fill(background)
	Draw(dst, a, b, src[0], foreground)

	// Result should be a filled circle
	for y := int32(0); y < dst.Height(); y++ {
		for x := int32(0); x < dst.Width(); x++ {
			r2 := (x-a)*(x-a) + (y-b)*(y-b)
			var e nimble.Pixel
			if r2 <= 25 {
				e = foreground
			} else {
				e = background
			}
			if e != dst.Pixel(x, y) {
				t.Fail()
				panic(fmt.Sprintf("x=%v y=%v e=%v\n", x, y, e))
			}
		}
	}
}

var src = MakeAnimation(10, false, 1)
var dst = nimble.MakePixMap(500, 500, make([]nimble.Pixel, 500*500), 500)

func BenchmarkDraw(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Draw(dst, 250, 250, src[0], nimble.Gray(1))
	}
}
