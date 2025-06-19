package ui

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// rangeSliderRenderer handles the layout and rendering of the RangeSlider widget.
type rangeSliderRenderer struct {
	fullTrackBar     *canvas.Rectangle   // The gray background bar that spans the slider
	selectedRangeBar *canvas.Rectangle   // The colored active range between start and end
	startThumb       fyne.CanvasObject   // The draggable thumb for the start of the range
	endThumb         fyne.CanvasObject   // The draggable thumb for the end of the range
	slider           *RangeSlider        // The main RangeSlider widget reference
	allObjects       []fyne.CanvasObject // All drawable objects returned in Objects()
}

// newRangeSliderRenderer initializes the visual elements of the range slider.
func newRangeSliderRenderer(s *RangeSlider) *rangeSliderRenderer {
	thumbSize := fyne.NewSize(14, 14)
	thumbColor := color.NRGBA{R: 100, G: 150, B: 255, A: 255}

	// Background bar representing the full range
	fullBar := canvas.NewRectangle(color.Gray{Y: 220})

	// Active bar showing the selected range
	activeBar := canvas.NewRectangle(thumbColor)

	// Thumbs for dragging start and end
	startThumb := NewPointerThumb(thumbColor, thumbSize)
	endThumb := NewPointerThumb(thumbColor, thumbSize)

	r := &rangeSliderRenderer{
		fullTrackBar:     fullBar,
		selectedRangeBar: activeBar,
		startThumb:       startThumb,
		endThumb:         endThumb,
		slider:           s,
		allObjects:       []fyne.CanvasObject{fullBar, activeBar, startThumb, endThumb},
	}

	// Store the renderer back in the slider and lay out initial positions
	s.renderer = r
	r.Layout(s.Size())
	return r
}

// Layout positions and sizes all elements in the range slider based on the current widget size and values.
func (r *rangeSliderRenderer) Layout(size fyne.Size) {
	slider := r.slider
	thumbSize := r.startThumb.MinSize()
	halfThumbWidth := thumbSize.Width / 2
	halfThumbHeight := thumbSize.Height / 2

	// Compute the usable horizontal space, excluding the thumb width
	usableWidth := float64(size.Width - thumbSize.Width)

	// Normalize a value to a 0–1 scale based on slider min and max
	percentOf := func(value float64) float64 {
		return (value - slider.Min) / (slider.Max - slider.Min)
	}

	// Compute X positions for the start and end thumbs
	startX := float32(percentOf(slider.Start)*usableWidth) + halfThumbWidth
	endX := float32(percentOf(slider.End)*usableWidth) + halfThumbWidth

	// Y-coordinate for the bars
	barY := size.Height/2 - 2

	// Position the full track bar
	r.fullTrackBar.Move(fyne.NewPos(halfThumbWidth, barY))
	r.fullTrackBar.Resize(fyne.NewSize(size.Width-thumbSize.Width, 4))

	// Calculate left and right edges of the active bar (between start and end)
	leftX := float32(math.Min(float64(startX), float64(endX)))
	rightX := float32(math.Max(float64(startX), float64(endX)))

	// Position and resize the selected (active) range bar
	r.selectedRangeBar.Move(fyne.NewPos(leftX, barY))
	r.selectedRangeBar.Resize(fyne.NewSize(rightX-leftX, 4))

	// Helper function to move a thumb to its correct position
	positionThumb := func(thumb fyne.CanvasObject, x float32) {
		thumb.Move(fyne.NewPos(x-halfThumbWidth, size.Height/2-halfThumbHeight))
		thumb.Resize(thumbSize)
	}
	positionThumb(r.startThumb, startX)
	positionThumb(r.endThumb, endX)
}

// MinSize defines the minimum space the slider needs to render correctly.
func (r *rangeSliderRenderer) MinSize() fyne.Size {
	return fyne.NewSize(624, 32)
}

// Refresh redraws the entire widget.
func (r *rangeSliderRenderer) Refresh() {
	canvas.Refresh(r.slider)
}

// BackgroundColor defines what color shows behind the widget (transparent in this case).
func (r *rangeSliderRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

// Objects returns all visual elements that make up the slider.
func (r *rangeSliderRenderer) Objects() []fyne.CanvasObject {
	return r.allObjects
}

// Destroy cleans up resources — not needed for this basic widget.
func (r *rangeSliderRenderer) Destroy() {}
