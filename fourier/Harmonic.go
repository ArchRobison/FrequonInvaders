package fourier

// Harmonic holds the information needed to draw the fourier view of a Frequon
type Harmonic struct {
	Ωx, Ωy    float32 // Angular velocities
	Phase     float32 // Phase at (0,0)
	Amplitude float32 // Amplitude
}
