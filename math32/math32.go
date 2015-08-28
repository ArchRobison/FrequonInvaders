// Package math32 implements math functions on float32
package math32

import (
	"math"
)

const Pi = math.Pi

// RoundToInt rounds a 32-bit float to the nearest int and returns it.
// If the rounded value does not fit in an int, the result is undefined.
func RoundToInt(x float32) int {
	return int(math.Floor(float64(x) + 0.5))
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Sincos(θ float32) (sin, cos float32) {
	y, x := math.Sincos(float64(θ))
	sin = float32(y)
	cos = float32(x)
	return
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Hypot(x float32, y float32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

func Atan2(y float32, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Max(a, b float32) float32 {
	if a < b {
		return b
	} else {
		return a
	}
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}
