// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type Canvas interface {
	Size() math.Size
	IsComplete() bool
	Complete()
	Push()
	Pop()
	AddClip(rect math.Rect)
	Clear(color Color)
	DrawCanvas(canvas Canvas, position math.Point)
	DrawTexture(texture Texture, bounds math.Rect)
	DrawRunes(font Font, runes []rune, points []math.Point, color Color)
	DrawLines(polygon Polygon, pen Pen)
	DrawPolygon(polygon Polygon, pen Pen, brush Brush)
	DrawRect(rect math.Rect, brush Brush)
	DrawRoundedRect(rect math.Rect, tl, tr, bl, br float32, p Pen, b Brush)
}
