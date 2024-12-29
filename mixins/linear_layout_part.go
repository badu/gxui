// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type LinearLayoutPart struct {
	outer               gxui.LinearLayoutOuter
	direction           gxui.Direction
	sizeMode            gxui.SizeMode
	horizontalAlignment gxui.HorizontalAlignment
	verticalAlignment   gxui.VerticalAlignment
}

func (l *LinearLayoutPart) Init(outer gxui.LinearLayoutOuter) {
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
		case gxui.Horizontal:
			switch l.verticalAlignment {
			case gxui.AlignTop:
				minor = cm.T
			case gxui.AlignMiddle:
				minor = (s.H - cs.H) / 2
			case gxui.AlignBottom:
				minor = s.H - cs.H
			}
		case gxui.Vertical:
			switch l.horizontalAlignment {
			case gxui.AlignLeft:
				minor = cm.L
			case gxui.AlignCenter:
				minor = (s.W - cs.W) / 2
			case gxui.AlignRight:
				minor = s.W - cs.W
			}
		}

		// Peform layout
		switch l.direction {
		case gxui.LeftToRight:
			major += cm.L
			child.Offset = math.Point{X: major, Y: minor}.Add(o)
			major += cs.W
			major += cm.R
			s.W -= cs.W + cm.W()
		case gxui.RightToLeft:
			major -= cm.R
			child.Offset = math.Point{X: major - cs.W, Y: minor}.Add(o)
			major -= cs.W
			major -= cm.L
			s.W -= cs.W + cm.W()
		case gxui.TopToBottom:
			major += cm.T
			child.Offset = math.Point{X: minor, Y: major}.Add(o)
			major += cs.H
			major += cm.B
			s.H -= cs.H + cm.H()
		case gxui.BottomToTop:
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

func (l *LinearLayoutPart) Direction() gxui.Direction {
	return l.direction
}

func (l *LinearLayoutPart) SetDirection(d gxui.Direction) {
	if l.direction != d {
		l.direction = d
		l.outer.Relayout()
	}
}

func (l *LinearLayoutPart) SizeMode() gxui.SizeMode {
	return l.sizeMode
}

func (l *LinearLayoutPart) SetSizeMode(mode gxui.SizeMode) {
	if l.sizeMode != mode {
		l.sizeMode = mode
		l.outer.Relayout()
	}
}

func (l *LinearLayoutPart) HorizontalAlignment() gxui.HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LinearLayoutPart) SetHorizontalAlignment(alignment gxui.HorizontalAlignment) {
	if l.horizontalAlignment != alignment {
		l.horizontalAlignment = alignment
		l.outer.Relayout()
	}
}

func (l *LinearLayoutPart) VerticalAlignment() gxui.VerticalAlignment {
	return l.verticalAlignment
}

func (l *LinearLayoutPart) SetVerticalAlignment(alignment gxui.VerticalAlignment) {
	if l.verticalAlignment != alignment {
		l.verticalAlignment = alignment
		l.outer.Relayout()
	}
}
