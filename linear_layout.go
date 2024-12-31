// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

// LinearLayout is a Container that lays out its child Controls into a column or
// row. The layout will always start by positioning the first (0'th) child, and
// then depending on the direction, will position each successive child either
// to the left, top, right or bottom of the preceding child Control.
// LinearLayout makes no effort to distribute remaining space evenly between the
// children - an child control that is laid out before others will reduce the
// remaining space given to the later children, even to the point that there is
// zero space remaining.
type LinearLayout interface {
	// LinearLayout extends the Control interface.
	Control

	// LinearLayout extends the Container interface.
	Container

	// Direction returns the direction of layout for this LinearLayout.
	Direction() Direction

	// Direction sets the direction of layout for this LinearLayout.
	SetDirection(Direction)

	// SizeMode returns the desired size behaviour for this LinearLayout.
	SizeMode() SizeMode

	// SetSizeMode sets the desired size behaviour for this LinearLayout.
	SetSizeMode(SizeMode)

	// HAlign returns the alignment of the child Controls when laying out TopToBottom or BottomToTop.
	// It has no effect when the layout direction is LeftToRight or RightToLeft.
	HorizontalAlignment() HAlign

	// SetHorizontalAlignment sets the alignment of the child Controls when laying out TopToBottom or BottomToTop.
	// It has no effect when the layout direction is LeftToRight or RightToLeft.
	SetHorizontalAlignment(HAlign)

	// VAlign returns the alignment of the child Controls when laying out LeftToRight or RightToLeft.
	// It has no effect when the layout direction is TopToBottom or BottomToTop.
	VerticalAlignment() VAlign

	// SetVerticalAlignment returns the alignment of the child Controls when laying out LeftToRight or RightToLeft.
	// It has no effect when the layout direction is TopToBottom or BottomToTop.
	SetVerticalAlignment(VAlign)

	// BorderPen returns the Pen used to draw the LinearLayout's border.
	BorderPen() Pen

	// SetBorderPen sets the Pen used to draw the LinearLayout's border.
	SetBorderPen(Pen)

	// BackgroundBrush returns the Brush used to fill the LinearLayout's background.
	BackgroundBrush() Brush

	// SetBackgroundBrush sets the Brush used to fill the LinearLayout's background.
	SetBackgroundBrush(Brush)
}

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
	ContainerBase
	LinearLayoutPart
	BackgroundBorderPainter
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
