// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/math"
)

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
	offset := l.parent.Padding().TopLeft()
	children := l.parent.Children()
	major := 0

	if l.direction.RightToLeft() || l.direction.BottomToTop() {
		if l.direction.RightToLeft() {
			major = size.Width
		} else {
			major = size.Height
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
				minor = childMargin.Top
			case AlignMiddle:
				minor = (size.Height - childSize.Height) / 2
			case AlignBottom:
				minor = size.Height - childSize.Height
			}
		case Vertical:
			switch l.horizontalAlignment {
			case AlignLeft:
				minor = childMargin.Left
			case AlignCenter:
				minor = (size.Width - childSize.Width) / 2
			case AlignRight:
				minor = size.Width - childSize.Width
			}
		}

		// Perform layout
		switch l.direction {
		case LeftToRight:
			major += childMargin.Left
			child.Offset = math.Point{X: major, Y: minor}.Add(offset)
			major += childSize.Width
			major += childMargin.Right
			size.Width -= childSize.Width + childMargin.Width()
		case RightToLeft:
			major -= childMargin.Right
			child.Offset = math.Point{X: major - childSize.Width, Y: minor}.Add(offset)
			major -= childSize.Width
			major -= childMargin.Left
			size.Width -= childSize.Width + childMargin.Width()
		case TopToBottom:
			major += childMargin.Top
			child.Offset = math.Point{X: minor, Y: major}.Add(offset)
			major += childSize.Height
			major += childMargin.Bottom
			size.Height -= childSize.Height + childMargin.Height()
		case BottomToTop:
			major -= childMargin.Bottom
			child.Offset = math.Point{X: minor, Y: major - childSize.Height}.Add(offset)
			major -= childSize.Height
			major -= childMargin.Top
			size.Height -= childSize.Height + childMargin.Height()
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
			offset.X += childBounds.Width()
		} else {
			offset.Y += childBounds.Height()
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

func (l *LinearLayoutImpl) Init(parent BaseContainerParent, canvasCreator CanvasCreator) {
	l.ContainerBase.Init(parent, canvasCreator)
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
