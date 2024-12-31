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
	parent          ControlBaseParent
	onDragStart     Event
	onDragEnd       Event
	backgroundColor Color
	foregroundColor Color
	styles          *StyleDefs
	isDragging      bool
}

func (b *SplitterBar) Init(parent ControlBaseParent, driver Driver, styles *StyleDefs) {
	b.ControlBase.Init(parent, driver)
	b.styles = styles
	b.parent = parent

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
	if b.onDragStart == nil {
		b.onDragStart = CreateEvent(func(MouseEvent) {})
	}

	return b.onDragStart.Listen(callback)
}

func (b *SplitterBar) OnDragEnd(callback func(event MouseEvent)) EventSubscription {
	if b.onDragEnd == nil {
		b.onDragEnd = CreateEvent(func(MouseEvent) {})
	}

	return b.onDragEnd.Listen(callback)
}

// parts.DrawPaintPart overrides
func (b *SplitterBar) Paint(canvas Canvas) {
	rect := b.parent.Size().Rect()
	canvas.DrawRect(rect, CreateBrush(b.backgroundColor))
	if b.foregroundColor != b.backgroundColor {
		canvas.DrawRect(rect.ContractI(1), CreateBrush(b.foregroundColor))
	}
}

// InputEventHandlerPart overrides
func (b *SplitterBar) MouseDown(event MouseEvent) {
	b.isDragging = true
	b.onDragStart.Emit(event)
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
		b.onDragEnd.Emit(we)
	})

	b.InputEventHandlerPart.MouseDown(event)
}
