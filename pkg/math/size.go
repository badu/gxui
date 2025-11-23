// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

type Size struct {
	Width, Height int
}

func (s Size) Point() Point {
	return Point{X: s.Width, Y: s.Height}
}

func (s Size) Vec2() Vec2 {
	return Vec2{X: float32(s.Width), Y: float32(s.Height)}
}

func (s Size) Rect() Rect {
	return CreateRect(0, 0, s.Width, s.Height)
}

func (s Size) CenteredRect() Rect {
	return CreateRect(-s.Width/2, -s.Height/2, s.Width/2, s.Height/2)
}

func (s Size) Scale(vec2 Vec2) Size {
	return Size{
		Width:  int(math.Ceil(float64(s.Width) * float64(vec2.X))),
		Height: int(math.Ceil(float64(s.Height) * float64(vec2.Y))),
	}
}
func (s Size) ScaleS(size float32) Size {
	return Size{
		Width:  int(math.Ceil(float64(s.Width) * float64(size))),
		Height: int(math.Ceil(float64(s.Height) * float64(size))),
	}
}

func (s Size) Expand(spacing Spacing) Size {
	return Size{Width: s.Width + spacing.Width(), Height: s.Height + spacing.Height()}
}

func (s Size) Contract(spacing Spacing) Size {
	return Size{Width: s.Width - spacing.Width(), Height: s.Height - spacing.Height()}
}

func (s Size) Add(size Size) Size {
	return Size{Width: s.Width + size.Width, Height: s.Height + size.Height}
}

func (s Size) Sub(size Size) Size {
	return Size{Width: s.Width - size.Width, Height: s.Height - size.Height}
}

func (s Size) Min(size Size) Size {
	return Size{Width: min(s.Width, size.Width), Height: min(s.Height, size.Height)}
}

func (s Size) Max(size Size) Size {
	return Size{Width: max(s.Width, size.Width), Height: max(s.Height, size.Height)}
}

func (s Size) Clamp(min, max Size) Size {
	return Size{Width: Clamp(s.Width, min.Width, max.Width), Height: Clamp(s.Height, min.Height, max.Height)}
}

func (s Size) WH() (w, h int) {
	return s.Width, s.Height
}

func (s Size) Area() int {
	return s.Width * s.Height
}

func (s Size) EdgeAlignedFit(outer Rect, edgePoint Point) Rect {
	r := s.CenteredRect().Offset(edgePoint).Constrain(outer)
	if topFits := edgePoint.Y+s.Height < outer.Max.Y; topFits {
		return r.OffsetY(edgePoint.Y - r.Min.Y)
	}
	if bottomFits := edgePoint.Y-s.Height >= outer.Min.Y; bottomFits {
		return r.OffsetY(edgePoint.Y - r.Max.Y)
	}
	if leftFits := edgePoint.X+s.Width < outer.Max.X; leftFits {
		return r.OffsetX(edgePoint.X - r.Min.X)
	}
	if rightFits := edgePoint.X-s.Width >= outer.Min.X; rightFits {
		return r.OffsetX(edgePoint.X - r.Max.X)
	}
	return r
}
