// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

type ScrollLayout interface {
	Control
	Parent
	SetChild(Control)
	Child() Control
	SetScrollAxis(horizontal, vertical bool)
	ScrollAxis() (horizontal, vertical bool)
	BorderPen() Pen
	SetBorderPen(Pen)
	BackgroundBrush() Brush
	SetBackgroundBrush(Brush)
}

type ScrollLayoutImpl struct {
	ContainerBase
	BackgroundBorderPainter
	parent       BaseContainerParent
	scrollOffset math.Point
	canScrollX   bool
	canScrollY   bool
	scrollBarX   *Child
	scrollBarY   *Child
	child        *Child
	innerSize    math.Size
}

func (l *ScrollLayoutImpl) Init(parent BaseContainerParent, driver Driver, styles *StyleDefs) {
	l.ContainerBase.Init(parent, driver)
	l.BackgroundBorderPainter.Init(parent)

	l.parent = parent
	l.canScrollX = true
	l.canScrollY = true
	scrollBarX := CreateScrollBar(driver, styles)
	scrollBarX.SetOrientation(Horizontal)
	scrollBarX.OnScroll(func(from, to int) { l.SetScrollOffset(math.Point{X: from, Y: l.scrollOffset.Y}) })
	scrollBarY := CreateScrollBar(driver, styles)
	scrollBarY.SetOrientation(Vertical)
	scrollBarY.OnScroll(func(from, to int) { l.SetScrollOffset(math.Point{X: l.scrollOffset.X, Y: from}) })
	l.scrollBarX = l.AddChild(scrollBarX)
	l.scrollBarY = l.AddChild(scrollBarY)
	l.SetMouseEventTarget(true)
}

func (l *ScrollLayoutImpl) LayoutChildren() {
	size := l.parent.Size().Contract(l.Padding())
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
		maxSize := l.innerSize
		if l.canScrollX {
			maxSize.W = math.MaxSize.W
		}
		if l.canScrollY {
			maxSize.H = math.MaxSize.H
		}
		childSize := l.child.Control.DesiredSize(math.ZeroSize, maxSize)
		l.child.Layout(childSize.Rect().Offset(l.scrollOffset.Neg()).Offset(offset))
		l.scrollBarX.Control.(ScrollBar).SetScrollLimit(childSize.W)
		l.scrollBarY.Control.(ScrollBar).SetScrollLimit(childSize.H)
	}

	l.SetScrollOffset(l.scrollOffset)
}

func (l *ScrollLayoutImpl) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (l *ScrollLayoutImpl) SetScrollOffset(scrollOffset math.Point) bool {
	var childSize math.Size
	if l.child != nil {
		childSize = l.child.Control.Size()
	}

	innerSize := l.innerSize
	scrollOffset = scrollOffset.Min(childSize.Sub(innerSize).Point()).Max(math.Point{})

	l.scrollBarX.Control.SetVisible(l.canScrollX && childSize.W > innerSize.W)
	l.scrollBarY.Control.SetVisible(l.canScrollY && childSize.H > innerSize.H)
	l.scrollBarX.Control.(ScrollBar).SetScrollPosition(l.scrollOffset.X, l.scrollOffset.X+innerSize.W)
	l.scrollBarY.Control.(ScrollBar).SetScrollPosition(l.scrollOffset.Y, l.scrollOffset.Y+innerSize.H)

	if l.scrollOffset != scrollOffset {
		l.scrollOffset = scrollOffset
		l.ReLayout()
		return true
	}

	return false
}

// InputEventHandlerPart override
func (l *ScrollLayoutImpl) MouseScroll(event MouseEvent) bool {
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
func (l *ScrollLayoutImpl) SetChild(control Control) {
	if l.child != nil {
		l.RemoveChild(l.child.Control)
	}
	if control != nil {
		l.child = l.AddChildAt(0, control)
	}
}

func (l *ScrollLayoutImpl) Child() Control {
	return l.child.Control
}

func (l *ScrollLayoutImpl) SetScrollAxis(horizontal, vertical bool) {
	if l.canScrollX != horizontal || l.canScrollY != vertical {
		l.canScrollX, l.canScrollY = horizontal, vertical
		l.scrollBarX.Control.SetVisible(horizontal)
		l.scrollBarY.Control.SetVisible(vertical)
		l.ReLayout()
	}
}

func (l *ScrollLayoutImpl) ScrollAxis() (horizontal, vertical bool) {
	return l.canScrollX, l.canScrollY
}
