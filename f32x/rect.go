package f32x

import (
	"image"

	"gioui.org/f32"
	"gioui.org/unit"
)

type Rectangle struct {
	Min, Max f32.Point
}

func (r *Rectangle) SwappedBounds() Rectangle {
	min, max := r.Min, r.Max
	if max.X < min.X {
		max.X = r.Min.X
		min.X = r.Max.X
	}
	if max.Y < min.Y {
		max.Y = r.Min.Y
		min.Y = r.Max.Y
	}
	return Rectangle{Min: min, Max: max}
}

func (r *Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r *Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r *Rectangle) ConvertToPixelspace(dp func(v unit.Dp) int) image.Rectangle {
	return image.Rect(dp(unit.Dp(r.Min.X)), dp(unit.Dp(r.Min.Y)), dp(unit.Dp(r.Max.X)), dp(unit.Dp(r.Max.Y)))
}
