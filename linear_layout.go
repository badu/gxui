// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

type LinearLayoutParent interface {
	Container
	Size() math.Size           // was outer.Sized
	SetSize(newSize math.Size) // was outer.Sized
}

type LinearLayoutPart struct {
	parent              LinearLayoutParent
	sizeMode            SizeMode
	direction           Direction
	horizontalAlignment HAlign
	verticalAlignment   VAlign
}

func (l *LinearLayoutPart) Init(parent LinearLayoutParent) {
	l.parent = parent
}

func (l *LinearLayoutPart) LayoutChildren() {
	size := l.parent.Size().Contract(l.parent.Padding())
	offset := l.parent.Padding().LT()
	children := l.parent.Children()
	major := 0

	if l.direction.RightToLeft() || l.direction.BottomToTop() {
		if l.direction.RightToLeft() {
			major = size.W
		} else {
			major = size.H
		}
	}

	for _, child := range children {
		childMargin := child.Control.Margin()
		childSize := child.Control.DesiredSize(math.ZeroSize, size.Contract(childMargin).Max(math.ZeroSize))
		child.Control.SetSize(childSize)

		// Calculate minor-axis alignment
		var minor int
		switch l.direction.Orientation() {
		case Horizontal:
			switch l.verticalAlignment {
			case AlignTop:
				minor = childMargin.T
			case AlignMiddle:
				minor = (size.H - childSize.H) / 2
			case AlignBottom:
				minor = size.H - childSize.H
			}
		case Vertical:
			switch l.horizontalAlignment {
			case AlignLeft:
				minor = childMargin.L
			case AlignCenter:
				minor = (size.W - childSize.W) / 2
			case AlignRight:
				minor = size.W - childSize.W
			}
		}

		// Perform layout
		switch l.direction {
		case LeftToRight:
			major += childMargin.L
			child.Offset = math.Point{X: major, Y: minor}.Add(offset)
			major += childSize.W
			major += childMargin.R
			size.W -= childSize.W + childMargin.W()
		case RightToLeft:
			major -= childMargin.R
			child.Offset = math.Point{X: major - childSize.W, Y: minor}.Add(offset)
			major -= childSize.W
			major -= childMargin.L
			size.W -= childSize.W + childMargin.W()
		case TopToBottom:
			major += childMargin.T
			child.Offset = math.Point{X: minor, Y: major}.Add(offset)
			major += childSize.H
			major += childMargin.B
			size.H -= childSize.H + childMargin.H()
		case BottomToTop:
			major -= childMargin.B
			child.Offset = math.Point{X: minor, Y: major - childSize.H}.Add(offset)
			major -= childSize.H
			major -= childMargin.T
			size.H -= childSize.H + childMargin.H()
		}
	}
}

func (l *LinearLayoutPart) DesiredSize(min, max math.Size) math.Size {
	if l.sizeMode.Fill() {
		return max
	}

	bounds := min.Rect()
	children := l.parent.Children()

	horizontal := l.direction.Orientation().Horizontal()
	offset := math.Point{X: 0, Y: 0}
	for _, child := range children {
		childSize := child.Control.DesiredSize(math.ZeroSize, max)
		childMargin := child.Control.Margin()
		childBounds := childSize.Expand(childMargin).Rect().Offset(offset)
		if horizontal {
			offset.X += childBounds.W()
		} else {
			offset.Y += childBounds.H()
		}
		bounds = bounds.Union(childBounds)
	}

	return bounds.Size().Expand(l.parent.Padding()).Clamp(min, max)
}

func (l *LinearLayoutPart) Direction() Direction {
	return l.direction
}

func (l *LinearLayoutPart) SetDirection(d Direction) {
	if l.direction == d {
		return
	}

	l.direction = d
	l.parent.ReLayout()
}

func (l *LinearLayoutPart) SizeMode() SizeMode {
	return l.sizeMode
}

func (l *LinearLayoutPart) SetSizeMode(mode SizeMode) {
	if l.sizeMode == mode {
		return
	}

	l.sizeMode = mode
	l.parent.ReLayout()
}

func (l *LinearLayoutPart) HorizontalAlignment() HAlign {
	return l.horizontalAlignment
}

func (l *LinearLayoutPart) SetHorizontalAlignment(alignment HAlign) {
	if l.horizontalAlignment == alignment {
		return
	}
	l.horizontalAlignment = alignment
	l.parent.ReLayout()
}

func (l *LinearLayoutPart) VerticalAlignment() VAlign {
	return l.verticalAlignment
}

func (l *LinearLayoutPart) SetVerticalAlignment(alignment VAlign) {
	if l.verticalAlignment == alignment {
		return
	}

	l.verticalAlignment = alignment
	l.parent.ReLayout()
}

type LinearLayoutImpl struct {
	LinearLayoutPart
	BackgroundBorderPainter
	ContainerBase
}

func (l *LinearLayoutImpl) Init(parent BaseContainerParent, driver Driver) {
	l.ContainerBase.Init(parent, driver)
	l.LinearLayoutPart.Init(parent)
	l.BackgroundBorderPainter.Init(parent)
	l.SetMouseEventTarget(true)
	l.SetBackgroundBrush(TransparentBrush)
	l.SetBorderPen(TransparentPen)
}

func (l *LinearLayoutImpl) Paint(canvas Canvas) {
	rect := l.Size().Rect()
	l.BackgroundBorderPainter.PaintBackground(canvas, rect)
	l.PaintChildrenPart.Paint(canvas)
	l.BackgroundBorderPainter.PaintBorder(canvas, rect)
}
