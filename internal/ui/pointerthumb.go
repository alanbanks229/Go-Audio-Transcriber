package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// PointerThumb is a circular thumb with a pointer cursor.
type PointerThumb struct {
	widget.BaseWidget
	fill color.Color
	size fyne.Size
}

func NewPointerThumb(fill color.Color, size fyne.Size) *PointerThumb {
	pt := &PointerThumb{fill: fill, size: size}
	pt.ExtendBaseWidget(pt)
	return pt
}

func (pt *PointerThumb) CreateRenderer() fyne.WidgetRenderer {
	circle := canvas.NewCircle(pt.fill)
	circle.StrokeColor = color.Black
	circle.StrokeWidth = 1
	circle.Resize(pt.size)
	return &pointerThumbRenderer{circle: circle, objects: []fyne.CanvasObject{circle}}
}

func (pt *PointerThumb) Cursor() desktop.Cursor {
	return desktop.PointerCursor
}

type pointerThumbRenderer struct {
	circle  *canvas.Circle
	objects []fyne.CanvasObject
}

func (r *pointerThumbRenderer) Layout(size fyne.Size)        { r.circle.Resize(size) }
func (r *pointerThumbRenderer) Refresh()                     { r.circle.Refresh() }
func (r *pointerThumbRenderer) MinSize() fyne.Size           { return fyne.NewSize(14, 14) }
func (r *pointerThumbRenderer) BackgroundColor() color.Color { return color.Transparent }
func (r *pointerThumbRenderer) Objects() []fyne.CanvasObject { return r.objects }
func (r *pointerThumbRenderer) Destroy()                     {}
