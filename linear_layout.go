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

	// HorizontalAlignment returns the alignment of the child Controls when laying
	// out TopToBottom or BottomToTop. It has no effect when the layout direction
	// is LeftToRight or RightToLeft.
	HorizontalAlignment() HorizontalAlignment

	// SetHorizontalAlignment sets the alignment of the child Controls when laying
	// out TopToBottom or BottomToTop. It has no effect when the layout direction
	// is LeftToRight or RightToLeft.
	SetHorizontalAlignment(HorizontalAlignment)

	// VerticalAlignment returns the alignment of the child Controls when laying
	// out LeftToRight or RightToLeft. It has no effect when the layout direction
	// is TopToBottom or BottomToTop.
	VerticalAlignment() VerticalAlignment

	// SetVerticalAlignment returns the alignment of the child Controls when
	// laying out LeftToRight or RightToLeft. It has no effect when the layout
	// direction is TopToBottom or BottomToTop.
	SetVerticalAlignment(VerticalAlignment)

	// BorderPen returns the Pen used to draw the LinearLayout's border.
	BorderPen() Pen

	// SetBorderPen sets the Pen used to draw the LinearLayout's border.
	SetBorderPen(Pen)

	// BackgroundBrush returns the Brush used to fill the LinearLayout's
	// background.
	BackgroundBrush() Brush

	// SetBackgroundBrush sets the Brush used to fill the LinearLayout's
	// background.
	SetBackgroundBrush(Brush)
}

type LinearLayoutParent interface {
	Container
	Size() math.Size           // was outer.Sized
	SetSize(newSize math.Size) // was outer.Sized
}

type LinearLayoutPart struct {
	parent              LinearLayoutParent
	direction           Direction
	sizeMode            SizeMode
	horizontalAlignment HorizontalAlignment
	verticalAlignment   VerticalAlignment
}

func (l *LinearLayoutPart) Init(parent LinearLayoutParent) {
	l.parent = parent
}

func (l *LinearLayoutPart) LayoutChildren() {
	s := l.parent.Size().Contract(l.parent.Padding())
	o := l.parent.Padding().LT()
	children := l.parent.Children()
	major := 0

	if l.direction.RightToLeft() || l.direction.BottomToTop() {
		if l.direction.RightToLeft() {
			major = s.W
		} else {
			major = s.H
		}
	}

	for _, child := range children {
		cm := child.Control.Margin()
		cs := child.Control.DesiredSize(math.ZeroSize, s.Contract(cm).Max(math.ZeroSize))
		child.Control.SetSize(cs)

		// Calculate minor-axis alignment
		var minor int
		switch l.direction.Orientation() {
		case Horizontal:
			switch l.verticalAlignment {
			case AlignTop:
				minor = cm.T
			case AlignMiddle:
				minor = (s.H - cs.H) / 2
			case AlignBottom:
				minor = s.H - cs.H
			}
		case Vertical:
			switch l.horizontalAlignment {
			case AlignLeft:
				minor = cm.L
			case AlignCenter:
				minor = (s.W - cs.W) / 2
			case AlignRight:
				minor = s.W - cs.W
			}
		}

		// Peform layout
		switch l.direction {
		case LeftToRight:
			major += cm.L
			child.Offset = math.Point{X: major, Y: minor}.Add(o)
			major += cs.W
			major += cm.R
			s.W -= cs.W + cm.W()
		case RightToLeft:
			major -= cm.R
			child.Offset = math.Point{X: major - cs.W, Y: minor}.Add(o)
			major -= cs.W
			major -= cm.L
			s.W -= cs.W + cm.W()
		case TopToBottom:
			major += cm.T
			child.Offset = math.Point{X: minor, Y: major}.Add(o)
			major += cs.H
			major += cm.B
			s.H -= cs.H + cm.H()
		case BottomToTop:
			major -= cm.B
			child.Offset = math.Point{X: minor, Y: major - cs.H}.Add(o)
			major -= cs.H
			major -= cm.T
			s.H -= cs.H + cm.H()
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
		cs := child.Control.DesiredSize(math.ZeroSize, max)
		cm := child.Control.Margin()
		cb := cs.Expand(cm).Rect().Offset(offset)
		if horizontal {
			offset.X += cb.W()
		} else {
			offset.Y += cb.H()
		}
		bounds = bounds.Union(cb)
	}

	return bounds.Size().Expand(l.parent.Padding()).Clamp(min, max)
}

func (l *LinearLayoutPart) Direction() Direction {
	return l.direction
}

func (l *LinearLayoutPart) SetDirection(d Direction) {
	if l.direction != d {
		l.direction = d
		l.parent.Relayout()
	}
}

func (l *LinearLayoutPart) SizeMode() SizeMode {
	return l.sizeMode
}

func (l *LinearLayoutPart) SetSizeMode(mode SizeMode) {
	if l.sizeMode != mode {
		l.sizeMode = mode
		l.parent.Relayout()
	}
}

func (l *LinearLayoutPart) HorizontalAlignment() HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LinearLayoutPart) SetHorizontalAlignment(alignment HorizontalAlignment) {
	if l.horizontalAlignment != alignment {
		l.horizontalAlignment = alignment
		l.parent.Relayout()
	}
}

func (l *LinearLayoutPart) VerticalAlignment() VerticalAlignment {
	return l.verticalAlignment
}

func (l *LinearLayoutPart) SetVerticalAlignment(alignment VerticalAlignment) {
	if l.verticalAlignment != alignment {
		l.verticalAlignment = alignment
		l.parent.Relayout()
	}
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
