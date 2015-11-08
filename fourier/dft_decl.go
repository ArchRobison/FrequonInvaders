// Declarations for DFT

package fourier

import (
	"github.com/ArchRobison/Gophetica/nimble"
)

// vector length
const vlen = 4

// vec is a vector of float
type vec [vlen]float32

// cvec is a complex vector of float
// Mathematically, it's isomorphic to a vector of complex,
// but enabled more efficient evaluation of complex arithmetic.
type cvec struct {
	re, im vec
}

// foot is a struct of vectors that holds sums for computing
// a + bi, (a+bi)*(c+di), and (a+bi)*(c-di)
//
// It's called foot because the computation pattern resembles
// a thee-toed foot:
//
//            *u3
//        z---------> (next foot)
//       / \
//      / | \
//    zu* z  zu
//
type foot struct {
	a, b, ac, bc, ad, bd vec
}

// pixelsPerFoot is the number of pixels computed from a foot
const pixelsPerFoot = 3 * vlen

// u13 holds a complex value u and its cube.
type u13 struct {
	u1, u3 complex64
}

func accumulateToFeet(z *[2]cvec, u *[2]u13, feet []foot)

func feetToPixel(feet []foot, clut *colorLookupTable, row []nimble.Pixel)

func rotate(w []cvec, v []complex64)
