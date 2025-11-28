package cgo

import (
	"fmt"

	"github.com/badu/gxui/pkg/math"
)

// 16:16 fixed point ratio of DIPs to pixels
type resolution uint32

func (r resolution) String() string {
	return fmt.Sprintf("%f", r.dipsToPixels())
}

func (r resolution) dipsToPixels() float32 {
	return float32(r) / 65536.0
}

func (r resolution) intDipsToPixels(size int) int {
	return (size * int(r)) >> 16
}

func (r resolution) pointDipsToPixels(point math.Point) math.Point {
	return math.Point{
		X: r.intDipsToPixels(point.X),
		Y: r.intDipsToPixels(point.Y),
	}
}

func (r resolution) sizeDipsToPixels(size math.Size) math.Size {
	return math.Size{
		Width:  r.intDipsToPixels(size.Width),
		Height: r.intDipsToPixels(size.Height),
	}
}

func (r resolution) rectDipsToPixels(rect math.Rect) math.Rect {
	return math.Rect{
		Min: r.pointDipsToPixels(rect.Min),
		Max: r.pointDipsToPixels(rect.Max),
	}
}
