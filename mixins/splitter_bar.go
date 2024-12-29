// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/mixins/base"
)

type SplitterBar struct {
	base.ControlBase
	onDrag          func(wndPnt math.Point)
	outer           base.ControlBaseOuter
	theme           gxui.Theme
	onDragStart     gxui.Event
	onDragEnd       gxui.Event
	backgroundColor gxui.Color
	foregroundColor gxui.Color
	isDragging      bool
}

func (b *SplitterBar) Init(outer base.ControlBaseOuter, theme gxui.Theme) {
	b.ControlBase.Init(outer, theme)

	b.outer = outer
	b.theme = theme
	b.onDragStart = gxui.CreateEvent(func(gxui.MouseEvent) {})
	b.onDragEnd = gxui.CreateEvent(func(gxui.MouseEvent) {})
	b.backgroundColor = gxui.Red
	b.foregroundColor = gxui.Green
}

func (b *SplitterBar) SetBackgroundColor(color gxui.Color) {
	b.backgroundColor = color
}

func (b *SplitterBar) SetForegroundColor(color gxui.Color) {
	b.foregroundColor = color
}

func (b *SplitterBar) OnSplitterDragged(callback func(point math.Point)) {
	b.onDrag = callback
}

func (b *SplitterBar) IsDragging() bool {
	return b.isDragging
}

func (b *SplitterBar) OnDragStart(callback func(event gxui.MouseEvent)) gxui.EventSubscription {
	return b.onDragStart.Listen(callback)
}

func (b *SplitterBar) OnDragEnd(callback func(event gxui.MouseEvent)) gxui.EventSubscription {
	return b.onDragEnd.Listen(callback)
}

// parts.DrawPaintPart overrides
func (b *SplitterBar) Paint(canvas gxui.Canvas) {
	rect := b.outer.Size().Rect()
	canvas.DrawRect(rect, gxui.CreateBrush(b.backgroundColor))
	if b.foregroundColor != b.backgroundColor {
		canvas.DrawRect(rect.ContractI(1), gxui.CreateBrush(b.foregroundColor))
	}
}

// InputEventHandlerPart overrides
func (b *SplitterBar) MouseDown(event gxui.MouseEvent) {
	b.isDragging = true
	b.onDragStart.Fire(event)
	var mms, mus gxui.EventSubscription
	mms = event.Window.OnMouseMove(func(we gxui.MouseEvent) {
		if b.onDrag != nil {
			b.onDrag(we.WindowPoint)
		}
	})
	mus = event.Window.OnMouseUp(func(we gxui.MouseEvent) {
		mms.Unlisten()
		mus.Unlisten()
		b.isDragging = false
		b.onDragEnd.Fire(we)
	})

	b.InputEventHandlerPart.MouseDown(event)
}
