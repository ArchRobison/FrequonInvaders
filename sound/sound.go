package sound

import (
	"github.com/ArchRobison/Gophetica/math32"
	"github.com/ArchRobison/Gophetica/nimble"
	"math/rand"
)

var Twang, AntiTwang, Broken, Wobble, Bell []float32

const π = math32.Pi

func makeSound(n int, f func(float32) float32) []float32 {
	wave := make([]float32, n)
	for i := range wave {
		wave[i] = f(float32(i))
	}
	return wave
}

func init() {
	Twang = makeSound(44100, func(i float32) float32 {
		const ω = 110 * 2 * π / nimble.SampleRate
		const nharmonic = 32
		sum := float32(0)
		for h := float32(1); h <= nharmonic; h++ {
			sum += math32.Sin(ω*h*i) * math32.Exp(-i*0.00004*h)
		}
		return sum / nharmonic
	})

	n := len(Twang)
	AntiTwang = make([]float32, n)
	for i := range AntiTwang {
		AntiTwang[i] = Twang[n-1-i]
	}

	sum := float32(0)
	Broken = makeSound(44100, func(i float32) float32 {
		r := (rand.Float32() - 0.5) * (1.0 / 16)
		newSum := sum + r
		if math32.Abs(newSum) < 1.0 {
			sum = newSum
		}
		sum *= (1 - 1/32.0)
		return sum * math32.Exp(-i*.0001)
	})

	Wobble = makeSound(44100, func(i float32) float32 {
		return 0.25 * math32.Sin((440+44*math32.Sin(i*0.001))*i*(2*π/nimble.SampleRate)) * math32.Exp(-i*0.00005)
	})

	partials := [...]float32{0.56, 0.92, 1.19, 1.71, 2.00, 2.74, 3.00, 3.76, 4.07}
	Bell = makeSound(44100, func(i float32) float32 {
		const θ = 512 * 2 * π / nimble.SampleRate
		sum := float32(0)
		const n = float32(len(partials))
		for j := range partials {
			sum += math32.Sin(i * θ * partials[j])
		}
		return sum * (1.0 / n) * math32.Exp(-i*0.0001)
	})
}

func Play(wave []float32, relativePitch float32) {
	nimble.PlaySound(wave, 1, relativePitch)
}
