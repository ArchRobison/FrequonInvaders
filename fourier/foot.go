package fourier

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

// vector length
const vlen = 4

// vector of float
type vec [vlen]float32

// complex vector of float
type cvec struct {
	re, im vec
}

// "foot" struct of vectors that holds sums for computing
// a + bi, (a+bi)*(c+di), and (a+bi)*(c-di)
type foot struct {
	a, b, ac, bc, ad, bd [vlen]float32
}

// Structure with w and w cubed.
type w13 struct {
	w1, w3 complex64
}

func accumulateToFeet(z *[2]cvec, u *[2]w13, feet []foot)

func feetToPixel(feet []foot, clut *[128][128]nimble.Pixel, row []nimble.Pixel)

func rotate([]cvec, []complex64)
