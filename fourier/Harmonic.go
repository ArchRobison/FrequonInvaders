package fourier

// Harmonic specifies the information to draw a fourier view of a Frequon
type Harmonic struct {
	Ωx, Ωy    float32 // Spatical frequencies in radians per pixel
	Phase     float32 // Phase at (0,0), i.e. upper left pixel
	Amplitude float32 // Amplitude
}
