// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"fmt"

	"github.com/badu/gxui/math"
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
		W: r.intDipsToPixels(size.W),
		H: r.intDipsToPixels(size.H),
	}
}

func (r resolution) rectDipsToPixels(rect math.Rect) math.Rect {
	return math.Rect{
		Min: r.pointDipsToPixels(rect.Min),
		Max: r.pointDipsToPixels(rect.Max),
	}
}
