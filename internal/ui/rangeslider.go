package ui

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type RangeSlider struct {
	widget.BaseWidget
	Min, Max   float64
	Start, End float64

	OnChanged                  func()
	draggingStart, draggingEnd bool
	renderer                   *rangeSliderRenderer
}

func NewRangeSlider(min, max, start, end float64) *RangeSlider {
	s := &RangeSlider{Min: min, Max: max, Start: start, End: end}
	s.ExtendBaseWidget(s)
	return s
}

var (
	_ fyne.Widget       = (*RangeSlider)(nil)
	_ fyne.Draggable    = (*RangeSlider)(nil)
	_ desktop.Mouseable = (*RangeSlider)(nil)
)

func (s *RangeSlider) CreateRenderer() fyne.WidgetRenderer {
	// helper: refer to rangeslider_renderer.go
	return newRangeSliderRenderer(s)
}

func (s *RangeSlider) SetBoundsAndValues(min, max, start, end float64) {
	s.Min = min
	s.Max = max
	s.Start = clamp(start, min, max)
	s.End = clamp(end, min, max)

	if s.renderer != nil {
		s.Refresh()
	}
}

func (s *RangeSlider) Dragged(e *fyne.DragEvent) {
	width := float64(s.Size().Width)
	if width == 0 {
		return
	}
	scale := func(px float32) float64 {
		return clamp(s.Min+(float64(px)/width)*(s.Max-s.Min), s.Min, s.Max)
	}

	if s.draggingStart {
		s.Start = min(scale(e.Position.X), s.End)
	} else if s.draggingEnd {
		s.End = max(scale(e.Position.X), s.Start)
	}

	if s.OnChanged != nil {
		s.OnChanged()
	}
}

func (s *RangeSlider) DragEnd() {
	s.draggingStart, s.draggingEnd = false, false
}

func (s *RangeSlider) MouseDown(e *desktop.MouseEvent) {
	checkHit := func(pos fyne.Position, thumb fyne.CanvasObject) bool {
		size := thumb.Size()
		return e.Position.X >= pos.X && e.Position.X <= pos.X+size.Width &&
			e.Position.Y >= pos.Y && e.Position.Y <= pos.Y+size.Height
	}

	if checkHit(s.renderer.startThumb.Position(), s.renderer.startThumb) {
		s.draggingStart = true
	} else if checkHit(s.renderer.endThumb.Position(), s.renderer.endThumb) {
		s.draggingEnd = true
	}
}

func (s *RangeSlider) MouseUp(_ *desktop.MouseEvent) {
	s.DragEnd()
}

// helpers
func clamp(val, min, max float64) float64 {
	return math.Max(min, math.Min(max, val))
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
