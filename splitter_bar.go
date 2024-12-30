// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type SplitterBar struct {
	ControlBase
	onDrag          func(wndPnt math.Point)
	outer           ControlBaseOuter
	theme           App
	onDragStart     Event
	onDragEnd       Event
	backgroundColor Color
	foregroundColor Color
	isDragging      bool
}

func (b *SplitterBar) Init(outer ControlBaseOuter, theme App) {
	b.ControlBase.Init(outer, theme)

	b.outer = outer
	b.theme = theme
	b.onDragStart = CreateEvent(func(MouseEvent) {})
	b.onDragEnd = CreateEvent(func(MouseEvent) {})
	b.backgroundColor = Red
	b.foregroundColor = Green
}

func (b *SplitterBar) SetBackgroundColor(color Color) {
	b.backgroundColor = color
}

func (b *SplitterBar) SetForegroundColor(color Color) {
	b.foregroundColor = color
}

func (b *SplitterBar) OnSplitterDragged(callback func(point math.Point)) {
	b.onDrag = callback
}

func (b *SplitterBar) IsDragging() bool {
	return b.isDragging
}

func (b *SplitterBar) OnDragStart(callback func(event MouseEvent)) EventSubscription {
	return b.onDragStart.Listen(callback)
}

func (b *SplitterBar) OnDragEnd(callback func(event MouseEvent)) EventSubscription {
	return b.onDragEnd.Listen(callback)
}

// parts.DrawPaintPart overrides
func (b *SplitterBar) Paint(canvas Canvas) {
	rect := b.outer.Size().Rect()
	canvas.DrawRect(rect, CreateBrush(b.backgroundColor))
	if b.foregroundColor != b.backgroundColor {
		canvas.DrawRect(rect.ContractI(1), CreateBrush(b.foregroundColor))
	}
}

// InputEventHandlerPart overrides
func (b *SplitterBar) MouseDown(event MouseEvent) {
	b.isDragging = true
	b.onDragStart.Fire(event)
	var mms, mus EventSubscription
	mms = event.Window.OnMouseMove(func(we MouseEvent) {
		if b.onDrag != nil {
			b.onDrag(we.WindowPoint)
		}
	})
	mus = event.Window.OnMouseUp(func(we MouseEvent) {
		mms.Forget()
		mus.Forget()
		b.isDragging = false
		b.onDragEnd.Fire(we)
	})

	b.InputEventHandlerPart.MouseDown(event)
}
