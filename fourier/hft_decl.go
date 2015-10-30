// Declarations for HFT

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
//
// It's called foot because the computation pattern resembles
// a thee-toed foot:
//
//		      *u3
//        z---------> (next foot)
//	     / \
//      / | \
//	  zu* z  zu
//
type foot struct {
	a, b, ac, bc, ad, bd vec
}

// Structure with u and u cubed.
type u13 struct {
	u1, u3 complex64
}

func accumulateToFeet(z *[2]cvec, u *[2]u13, feet []foot)

func feetToPixel(feet []foot, clut *[128][128]nimble.Pixel, row []nimble.Pixel)

func rotate(w []cvec, v []complex64)
