package audio

// AudioSystem defines the interface for audio systems
type AudioSystem interface {
	Start() error
	Stop()
	SetBeepState(playing bool)
	Close()
}

// Config holds audio configuration
type Config struct {
	Frequency     int     // Sample rate (e.g., 44100)
	BeepFrequency float64 // Frequency of the beep sound (e.g., 440.0)
	Volume        float32 // Volume level (0.0 to 1.0)
}
