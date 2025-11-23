// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

type BackgroundBorderParent interface {
	Redraw() // was outer.Redrawer
}

type BackgroundBorderPainter struct {
	parent BackgroundBorderParent
	brush  Brush
	pen    Pen
}

func (b *BackgroundBorderPainter) Init(parent BackgroundBorderParent) {
	b.parent = parent
	b.brush = DefaultBrush
	b.pen = DefaultPen
}

func (b *BackgroundBorderPainter) PaintBackground(canvas Canvas, rect math.Rect) {
	if b.brush.Color.A == 0 {
		return
	}

	w := b.pen.Width
	canvas.DrawRoundedRect(rect, w, w, w, w, TransparentPen, b.brush)
}

func (b *BackgroundBorderPainter) PaintBorder(canvas Canvas, rect math.Rect) {
	if b.pen.Color.A != 0 && b.pen.Width != 0 {
		w := b.pen.Width
		canvas.DrawRoundedRect(rect, w, w, w, w, b.pen, TransparentBrush)
	}
}

func (b *BackgroundBorderPainter) BackgroundBrush() Brush {
	return b.brush
}

func (b *BackgroundBorderPainter) SetBackgroundBrush(brush Brush) {
	if b.brush == brush {
		return
	}

	b.brush = brush
	b.parent.Redraw()
}

func (b *BackgroundBorderPainter) BorderPen() Pen {
	return b.pen
}

func (b *BackgroundBorderPainter) SetBorderPen(pen Pen) {
	if b.pen == pen {
		return
	}

	b.pen = pen
	b.parent.Redraw()
}
