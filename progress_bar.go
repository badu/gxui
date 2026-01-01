// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

type ProgressBarParent interface {
	ControlBaseParent
	PaintProgress(Canvas, math.Rect, float32)
}

type ProgressBarImpl struct {
	InputEventHandlerPart
	ParentablePart
	DrawPaintPart
	AttachablePart
	VisiblePart
	LayoutablePart
	BackgroundBorderPainter
	parent           ProgressBarParent
	desiredSize      math.Size
	progress, target int
}

func (b *ProgressBarImpl) Init(parent ProgressBarParent, driver Driver, styles *StyleDefs) {
	b.parent = parent
	b.DrawPaintPart.Init(parent, driver)
	b.LayoutablePart.Init(parent)
	b.InputEventHandlerPart.Init()
	b.VisiblePart.Init(parent)
	b.BackgroundBorderPainter.Init(parent)
	b.desiredSize = math.MaxSize
	b.target = 100
}

func (b *ProgressBarImpl) Paint(canvas Canvas) {
	fraction := math.Saturate(float32(b.progress) / float32(b.target))
	rect := b.parent.Size().Rect()
	b.PaintBackground(canvas, rect)
	b.parent.PaintProgress(canvas, rect, fraction)
	b.PaintBorder(canvas, rect)
}

func (b *ProgressBarImpl) PaintProgress(canvas Canvas, rect math.Rect, fraction float32) {
	rect.Max.X = math.Lerp(rect.Min.X, rect.Max.X, fraction)
	canvas.DrawRect(rect, CreateBrush(Gray50))
}

func (b *ProgressBarImpl) DesiredSize(min, max math.Size) math.Size {
	return b.desiredSize.Clamp(min, max)
}

func (b *ProgressBarImpl) SetDesiredSize(size math.Size) {
	b.desiredSize = size
	b.ReLayout()
}

func (b *ProgressBarImpl) SetProgress(progress int) {
	if b.progress != progress {
		b.progress = progress
		b.Redraw()
	}
}

func (b *ProgressBarImpl) Progress() int {
	return b.progress
}

func (b *ProgressBarImpl) SetTarget(target int) {
	if b.target != target {
		b.target = target
		b.Redraw()
	}
}

func (b *ProgressBarImpl) Target() int {
	return b.target
}
