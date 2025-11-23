// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Spacing struct {
	Left, Top, Right, Bottom int
}

func CreateSpacing(size int) Spacing {
	return Spacing{Left: size, Top: size, Right: size, Bottom: size}
}

func (s Spacing) TopLeft() Point {
	return Point{X: s.Left, Y: s.Top}
}

func (s Spacing) Width() int {
	return s.Left + s.Right
}

func (s Spacing) Height() int {
	return s.Top + s.Bottom
}

func (s Spacing) Size() Size {
	return Size{Width: s.Width(), Height: s.Height()}
}

func (s Spacing) Add(spacing Spacing) Spacing {
	return Spacing{Left: s.Left + spacing.Left, Top: s.Top + spacing.Top, Right: s.Right + spacing.Right, Bottom: s.Bottom + spacing.Bottom}
}

func (s Spacing) Sub(spacing Spacing) Spacing {
	return Spacing{Left: s.Left - spacing.Left, Top: s.Top - spacing.Top, Right: s.Right - spacing.Right, Bottom: s.Bottom - spacing.Bottom}
}

func (s Spacing) Min(spacing Spacing) Spacing {
	return Spacing{Left: min(s.Left, spacing.Left), Top: min(s.Top, spacing.Top), Right: min(s.Right, spacing.Right), Bottom: min(s.Bottom, spacing.Bottom)}
}

func (s Spacing) Max(spacing Spacing) Spacing {
	return Spacing{Left: max(s.Left, spacing.Left), Top: max(s.Top, spacing.Top), Right: max(s.Right, spacing.Right), Bottom: max(s.Bottom, spacing.Bottom)}
}
