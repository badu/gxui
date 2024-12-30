// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type LinearLayoutPart struct {
	outer               LinearLayoutOuter
	direction           Direction
	sizeMode            SizeMode
	horizontalAlignment HorizontalAlignment
	verticalAlignment   VerticalAlignment
}

func (l *LinearLayoutPart) Init(outer LinearLayoutOuter) {
	l.outer = outer
}

func (l *LinearLayoutPart) LayoutChildren() {
	s := l.outer.Size().Contract(l.outer.Padding())
	o := l.outer.Padding().LT()
	children := l.outer.Children()
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
	children := l.outer.Children()

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

	return bounds.Size().Expand(l.outer.Padding()).Clamp(min, max)
}

func (l *LinearLayoutPart) Direction() Direction {
	return l.direction
}

func (l *LinearLayoutPart) SetDirection(d Direction) {
	if l.direction != d {
		l.direction = d
		l.outer.Relayout()
	}
}

func (l *LinearLayoutPart) SizeMode() SizeMode {
	return l.sizeMode
}

func (l *LinearLayoutPart) SetSizeMode(mode SizeMode) {
	if l.sizeMode != mode {
		l.sizeMode = mode
		l.outer.Relayout()
	}
}

func (l *LinearLayoutPart) HorizontalAlignment() HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LinearLayoutPart) SetHorizontalAlignment(alignment HorizontalAlignment) {
	if l.horizontalAlignment != alignment {
		l.horizontalAlignment = alignment
		l.outer.Relayout()
	}
}

func (l *LinearLayoutPart) VerticalAlignment() VerticalAlignment {
	return l.verticalAlignment
}

func (l *LinearLayoutPart) SetVerticalAlignment(alignment VerticalAlignment) {
	if l.verticalAlignment != alignment {
		l.verticalAlignment = alignment
		l.outer.Relayout()
	}
}
