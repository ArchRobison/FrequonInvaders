package fourier

// vector length
const vlen = 4

type vec [vlen]float32

// complex vector
type cvec struct {
	re, im vec
}

func (z *cvec) put(i int, val complex128) {
	z.re[i] = float32(real(val))
	z.im[i] = float32(imag(val))
}

func (z *cvec) get(i int) complex128 {
	return complex(float64(z.re[i]), float64(z.im[i]))
}

func put(z []cvec, i int, val complex128) {
	z[i/4].put(i%4, val)
}

func get(z []cvec, i int) complex128 {
	return z[i/4].get(i % 4)
}

// SOA "foot" structure holds sums for computing 
// a + bi, (a+bi)*(c+di), and (a+bi)*(c-di)
type foot struct {
	a, b, ac, bc, ad, bd [4]float32
}

type w13 struct {
	w1, w3 complex64
}

func accumulateToFeet(z *[2]cvec, u *[2]w13, feet []foot)

func accumulateToFeetSlow(z *[2]cvec, u *[2]w13, feet []foot) {
	for i := range feet {
		f := &feet[i]
		for j := 0; j < 2; j++ {
			a := &(*z)[j]
			for k := 0; k < vlen; k++ {
				w1 := (*u)[j].w1
				ar := a.re[k]
				ai := a.im[k]
				f.a[k] += ar
				f.b[k] += ai
				f.ac[k] += ar * real(w1)
				f.bc[k] += ai * real(w1)
				f.ad[k] += ar * imag(w1)
				f.bd[k] += ai * imag(w1)
				t := complex(ar, ai) * (*u)[j].w3
				a.re[k] = real(t)
				a.im[k] = imag(t)
			}
		}
	}
}

type pixel uint32

func feetToPixel(feet []foot, clut *[128][128]pixel, row []pixel)

func feetToPixelSlow(feet []foot, clut *[128][128]pixel, row []pixel) {
	for i := range feet {
		f := &feet[i]
		for k := 0; k < vlen; k++ {
			row[(3*i+0)*vlen+k] = clut[int(f.bc[k]-f.ad[k])][int(f.ac[k]+f.bd[k])]
			row[(3*i+1)*vlen+k] = clut[int(f.b[k])][int(f.a[k])]
			row[(3*i+2)*vlen+k] = clut[int(f.bc[k]+f.ad[k])][int(f.ac[k]-f.bd[k])]
			f.a[k] = 64.5
			f.b[k] = 64.5
			f.ac[k] = 64.5
			f.bc[k] = 64.5
			f.ad[k] = 0
		    f.bd[k] = 0
		}
	}

}

func rotate([]cvec, []complex64)

func rotateSlow(z []cvec, v []complex64) {
	for i := range z {
		for k := 0; k < vlen; k++ {
			t := complex(z[i].re[k], z[i].im[k]) * v[i]
			z[i].re[k] = real(t)
			z[i].im[k] = imag(t)
		}
	}
}
