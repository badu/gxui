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
