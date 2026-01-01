// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

type SplitterBarParent interface {
	Attached() bool
	Parent() Parent
	Size() math.Size
	Paint(canvas Canvas)
	Redraw()
}

type SplitterBar struct {
	InputEventHandlerPart
	ParentablePart
	DrawPaintPart
	AttachablePart
	VisiblePart
	LayoutablePart
	parent          SplitterBarParent
	onDragStart     Event
	onDragEnd       Event
	onDrag          func(point math.Point)
	styles          *StyleDefs
	BackgroundColor Color
	ForegroundColor Color
	IsDragging      bool
}

func (b *SplitterBar) Init(controlBaseParent SplitterBarParent, canvasCreator CanvasCreator, styles *StyleDefs) {
	b.DrawPaintPart.Init(controlBaseParent, canvasCreator)
	b.LayoutablePart.Init(controlBaseParent)
	b.InputEventHandlerPart.Init()
	b.VisiblePart.Init(controlBaseParent)

	b.styles = styles
	b.parent = controlBaseParent

	b.BackgroundColor = Red
	b.ForegroundColor = Green
}

func (b *SplitterBar) OnSplitterDragged(callback func(point math.Point)) {
	b.onDrag = callback
}

func (b *SplitterBar) OnDragStart(callback func(event MouseEvent)) EventSubscription {
	if b.onDragStart == nil {
		b.onDragStart = NewListener(func(MouseEvent) {})
	}

	return b.onDragStart.Listen(callback)
}

func (b *SplitterBar) OnDragEnd(callback func(event MouseEvent)) EventSubscription {
	if b.onDragEnd == nil {
		b.onDragEnd = NewListener(func(MouseEvent) {})
	}

	return b.onDragEnd.Listen(callback)
}

// parts.DrawPaintPart overrides
func (b *SplitterBar) Paint(canvas Canvas) {
	rect := b.parent.Size().Rect()
	canvas.DrawRect(rect, CreateBrush(b.BackgroundColor))
	if b.ForegroundColor != b.BackgroundColor {
		canvas.DrawRect(rect.ContractI(1), CreateBrush(b.ForegroundColor))
	}
}

// InputEventHandlerPart overrides
func (b *SplitterBar) MouseDown(event MouseEvent) {
	b.IsDragging = true
	b.onDragStart.Emit(event)
	var mms, mus EventSubscription
	mms = event.Window.OnMouseMove(
		func(windowEvent MouseEvent) {
			if b.onDrag == nil {
				return
			}

			b.onDrag(windowEvent.WindowPoint)
		},
	)
	mus = event.Window.OnMouseUp(
		func(windowEvent MouseEvent) {
			mms.Forget()
			mus.Forget()
			b.IsDragging = false
			b.onDragEnd.Emit(windowEvent)
		},
	)

	b.InputEventHandlerPart.MouseDown(event)
}

func (b *SplitterBar) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (b *SplitterBar) ContainsPoint(point math.Point) bool {
	return b.IsVisible() && b.Size().Rect().Contains(point)
}
