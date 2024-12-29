// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type ScrollLayout struct {
	ContainerBase
	BackgroundBorderPainter
	outer                  gxui.ContainerBaseOuter
	theme                  gxui.Theme
	scrollOffset           math.Point
	canScrollX, canScrollY bool
	scrollBarX, scrollBarY *gxui.Child
	child                  *gxui.Child
	innerSize              math.Size
}

func (l *ScrollLayout) Init(outer gxui.ContainerBaseOuter, theme gxui.Theme) {
	l.ContainerBase.Init(outer, theme)
	l.BackgroundBorderPainter.Init(outer)

	l.outer = outer
	l.theme = theme
	l.canScrollX = true
	l.canScrollY = true
	scrollBarX := theme.CreateScrollBar()
	scrollBarX.SetOrientation(gxui.Horizontal)
	scrollBarX.OnScroll(func(from, to int) { l.SetScrollOffset(math.Point{X: from, Y: l.scrollOffset.Y}) })
	scrollBarY := theme.CreateScrollBar()
	scrollBarY.SetOrientation(gxui.Vertical)
	scrollBarY.OnScroll(func(from, to int) { l.SetScrollOffset(math.Point{X: l.scrollOffset.X, Y: from}) })
	l.scrollBarX = l.AddChild(scrollBarX)
	l.scrollBarY = l.AddChild(scrollBarY)
	l.SetMouseEventTarget(true)
}

func (l *ScrollLayout) LayoutChildren() {
	size := l.outer.Size().Contract(l.Padding())
	offset := l.Padding().LT()

	var sxs, sys math.Size
	if l.canScrollX {
		sxs = l.scrollBarX.Control.DesiredSize(math.ZeroSize, size)
	}
	if l.canScrollY {
		sys = l.scrollBarY.Control.DesiredSize(math.ZeroSize, size)
	}

	l.scrollBarX.Layout(math.CreateRect(0, size.H-sxs.H, size.W-sys.W, size.H).Canon().Offset(offset))
	l.scrollBarY.Layout(math.CreateRect(size.W-sys.W, 0, size.W, size.H-sxs.H).Canon().Offset(offset))

	l.innerSize = size.Contract(math.Spacing{R: sys.W, B: sxs.H})

	if l.child != nil {
		max := l.innerSize
		if l.canScrollX {
			max.W = math.MaxSize.W
		}
		if l.canScrollY {
			max.H = math.MaxSize.H
		}
		cs := l.child.Control.DesiredSize(math.ZeroSize, max)
		l.child.Layout(cs.Rect().Offset(l.scrollOffset.Neg()).Offset(offset))
		l.scrollBarX.Control.(gxui.ScrollBar).SetScrollLimit(cs.W)
		l.scrollBarY.Control.(gxui.ScrollBar).SetScrollLimit(cs.H)
	}

	l.SetScrollOffset(l.scrollOffset)
}

func (l *ScrollLayout) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *ScrollLayout) SetScrollOffset(scrollOffset math.Point) bool {
	var childSize math.Size
	if l.child != nil {
		childSize = l.child.Control.Size()
	}

	innerSize := l.innerSize
	scrollOffset = scrollOffset.Min(childSize.Sub(innerSize).Point()).Max(math.Point{})

	l.scrollBarX.Control.SetVisible(l.canScrollX && childSize.W > innerSize.W)
	l.scrollBarY.Control.SetVisible(l.canScrollY && childSize.H > innerSize.H)
	l.scrollBarX.Control.(gxui.ScrollBar).SetScrollPosition(l.scrollOffset.X, l.scrollOffset.X+innerSize.W)
	l.scrollBarY.Control.(gxui.ScrollBar).SetScrollPosition(l.scrollOffset.Y, l.scrollOffset.Y+innerSize.H)

	if l.scrollOffset != scrollOffset {
		l.scrollOffset = scrollOffset
		l.Relayout()
		return true
	}

	return false
}

// InputEventHandlerPart override
func (l *ScrollLayout) MouseScroll(event gxui.MouseEvent) bool {
	if event.ScrollY == 0 {
		return l.InputEventHandlerPart.MouseScroll(event)
	}

	switch {
	case l.canScrollY:
		return l.SetScrollOffset(l.scrollOffset.AddY(-event.ScrollY))
	case l.canScrollX:
		return l.SetScrollOffset(l.scrollOffset.AddX(-event.ScrollY))
	default:
		return false
	}
}

// gxui.ScrollLayout complaince
func (l *ScrollLayout) SetChild(control gxui.Control) {
	if l.child != nil {
		l.RemoveChild(l.child.Control)
	}
	if control != nil {
		l.child = l.AddChildAt(0, control)
	}
}

func (l *ScrollLayout) Child() gxui.Control {
	return l.child.Control
}

func (l *ScrollLayout) SetScrollAxis(horizontal, vertical bool) {
	if l.canScrollX != horizontal || l.canScrollY != vertical {
		l.canScrollX, l.canScrollY = horizontal, vertical
		l.scrollBarX.Control.SetVisible(horizontal)
		l.scrollBarY.Control.SetVisible(vertical)
		l.Relayout()
	}
}

func (l *ScrollLayout) ScrollAxis() (horizontal, vertical bool) {
	return l.canScrollX, l.canScrollY
}
