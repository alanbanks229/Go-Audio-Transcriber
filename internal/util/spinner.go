package util

import (
	"fmt"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Spinner struct {
	frames     []*canvas.Image
	current    int
	ticker     *time.Ticker
	wrap       *fyne.Container
	isSpinning bool
}

func NewSpinner() *Spinner {
	frames := make([]*canvas.Image, 12)

	for i := 0; i < 12; i++ {
		path := ResolveAssetPath(filepath.Join("spinner_frames", fmt.Sprintf("spinner_%d.svg", i)))
		img := canvas.NewImageFromFile(path)
		img.SetMinSize(fyne.NewSize(32, 32))
		img.Hide()
		frames[i] = img
	}

	return &Spinner{
		frames: frames,
		wrap:   container.NewMax(frames[0]), // will swap in Start()
	}
}

// Widget returns the CanvasObject for layout (place this in your UI)
func (s *Spinner) Widget() fyne.CanvasObject {
	return s.wrap
}

// Start begins the animation loop
func (s *Spinner) Start() {
	if s.isSpinning {
		fmt.Println("Already Spinning")
		return
	}
	s.isSpinning = true

	s.ticker = time.NewTicker(83 * time.Millisecond) // ~12fps

	go func() {
		for range s.ticker.C {
			fyne.Do(func() {
				// Swap frame in container
				s.wrap.Objects = []fyne.CanvasObject{s.frames[s.current]}
				s.frames[s.current].Show()
				s.frames[s.current].Refresh()
				s.wrap.Refresh()
				s.current = (s.current + 1) % len(s.frames)
			})
		}
	}()
}

// Stop ends the animation and hides the spinner
func (s *Spinner) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.isSpinning = false
	s.current = 0

	fyne.Do(func() {
		for _, frame := range s.frames {
			frame.Hide()
		}
		s.wrap.Objects = []fyne.CanvasObject{}
		s.wrap.Refresh()
	})
}
