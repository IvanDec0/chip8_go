package audio

import (
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	defaultSamples = 512 //
)

// SDLAudioSystem implements AudioSystem using SDL2
type SDLAudioSystem struct {
	device    sdl.AudioDeviceID
	spec      sdl.AudioSpec
	phase     float64
	isPlaying bool
	buffer    []float32
	config    Config
}

// NewSDLAudioSystem creates a new SDL-based audio system
func NewSDLAudioSystem(config Config) (*SDLAudioSystem, error) {
	spec := &sdl.AudioSpec{
		Freq:     int32(config.Frequency),
		Format:   sdl.AUDIO_F32SYS,
		Channels: 1,
		Samples:  defaultSamples,
	}

	device, err := sdl.OpenAudioDevice("", false, spec, nil, sdl.AUDIO_ALLOW_ANY_CHANGE)
	if err != nil {
		return nil, err
	}

	return &SDLAudioSystem{
		device: device,
		spec:   *spec,
		buffer: make([]float32, defaultSamples),
		config: config,
	}, nil
}

// Start initializes the audio system
func (s *SDLAudioSystem) Start() error {
	return nil // SDL device is already initialized
}

// Stop stops the audio system
func (s *SDLAudioSystem) Stop() {
	s.SetBeepState(false)
}

// SetBeepState controls whether the beep sound should play
func (s *SDLAudioSystem) SetBeepState(playing bool) {
	if playing && !s.isPlaying {
		s.isPlaying = true
		s.generateAndQueueAudio()
		sdl.PauseAudioDevice(s.device, false)
	} else if !playing && s.isPlaying {
		s.isPlaying = false
		sdl.PauseAudioDevice(s.device, true)
		sdl.ClearQueuedAudio(s.device)
	}
}

// Update should be called regularly to maintain audio playback
func (s *SDLAudioSystem) Update() {
	if s.isPlaying {
		queuedSize := sdl.GetQueuedAudioSize(s.device)
		if queuedSize < uint32(len(s.buffer)*4*2) {
			s.generateAndQueueAudio()
		}
	}
}

// Close shuts down the audio system
func (s *SDLAudioSystem) Close() {
	if s.device != 0 {
		sdl.CloseAudioDevice(s.device)
	}
}

// generateAndQueueAudio generates sine wave audio and queues it
func (s *SDLAudioSystem) generateAndQueueAudio() {
	if !s.isPlaying {
		return
	}

	// Generate sine wave samples
	for i := range s.buffer {
		sample := float32(math.Sin(s.phase * 2 * math.Pi))
		s.buffer[i] = sample * s.config.Volume

		s.phase += s.config.BeepFrequency / float64(s.spec.Freq)
		if s.phase >= 1.0 {
			s.phase -= 1.0
		}
	}

	// Convert float32 slice to byte slice for SDL
	byteBuffer := (*[defaultSamples * 4]byte)(unsafe.Pointer(&s.buffer[0]))[: defaultSamples*4 : defaultSamples*4]
	sdl.QueueAudio(s.device, byteBuffer)
}
