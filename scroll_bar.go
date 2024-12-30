// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import "github.com/badu/gxui/math"

type ScrollBar interface {
	Control
	OnScroll(func(from, to int)) EventSubscription
	ScrollPosition() (from, to int)
	SetScrollPosition(from, to int)
	ScrollLimit() int
	SetScrollLimit(l int)
	AutoHide() bool
	SetAutoHide(l bool)
	Orientation() Orientation
	SetOrientation(Orientation)
}

type ScrollBarImpl struct {
	ControlBase
	parent             ControlBaseParent
	orientation        Orientation
	thickness          int
	minBarLength       int
	scrollPositionFrom int
	scrollPositionTo   int
	scrollLimit        int
	railBrush          Brush
	barBrush           Brush
	railPen            Pen
	barPen             Pen
	barRect            math.Rect
	onScroll           Event
	autoHide           bool
}

func (s *ScrollBarImpl) positionAt(p math.Point) int {
	orientation := s.orientation
	fraction := float32(orientation.Major(p.XY())) / float32(orientation.Major(s.Size().WH()))
	limit := s.ScrollLimit()
	return int(float32(limit) * fraction)
}

func (s *ScrollBarImpl) rangeAt(point math.Point) (int, int) {
	width := s.scrollPositionTo - s.scrollPositionFrom
	from := math.Clamp(s.positionAt(point), 0, s.scrollLimit-width)
	to := from + width
	return from, to
}

func (s *ScrollBarImpl) updateBarRect() {
	fractionFrom, fractionTo := s.ScrollFraction()
	size := s.Size()
	rect := size.Rect()
	halfMinLen := s.minBarLength / 2
	if s.orientation.Horizontal() {
		rect.Min.X = math.Lerp(0, size.W, fractionFrom)
		rect.Max.X = math.Lerp(0, size.W, fractionTo)
		if rect.W() < s.minBarLength {
			half := (rect.Min.X + rect.Max.X) / 2
			half = math.Clamp(half, rect.Min.X+halfMinLen, rect.Max.X-halfMinLen)
			rect.Min.X, rect.Max.X = half-halfMinLen, half+halfMinLen
		}
	} else {
		rect.Min.Y = math.Lerp(0, size.H, fractionFrom)
		rect.Max.Y = math.Lerp(0, size.H, fractionTo)
		if rect.H() < s.minBarLength {
			half := (rect.Min.Y + rect.Max.Y) / 2
			half = math.Clamp(half, rect.Min.Y+halfMinLen, rect.Max.Y-halfMinLen)
			rect.Min.Y, rect.Max.Y = half-halfMinLen, half+halfMinLen
		}
	}
	s.barRect = rect
}

func (s *ScrollBarImpl) Init(parent ControlBaseParent, driver Driver) {
	s.ControlBase.Init(parent, driver)

	s.parent = parent
	s.thickness = 10
	s.minBarLength = 10
	s.scrollPositionFrom = 0
	s.scrollPositionTo = 100
	s.scrollLimit = 100
	s.onScroll = CreateEvent(s.SetScrollPosition)
}

func (s *ScrollBarImpl) OnScroll(callback func(from, to int)) EventSubscription {
	return s.onScroll.Listen(callback)
}

func (s *ScrollBarImpl) ScrollFraction() (float32, float32) {
	from := float32(s.scrollPositionFrom) / float32(s.scrollLimit)
	to := float32(s.scrollPositionTo) / float32(s.scrollLimit)
	return from, to
}

func (s *ScrollBarImpl) DesiredSize(min, max math.Size) math.Size {
	if s.orientation.Horizontal() {
		return math.Size{W: max.W, H: s.thickness}.Clamp(min, max)
	} else {
		return math.Size{W: s.thickness, H: max.H}.Clamp(min, max)
	}
}

func (s *ScrollBarImpl) Paint(canvas Canvas) {
	canvas.DrawRoundedRect(s.parent.Size().Rect(), 3, 3, 3, 3, s.railPen, s.railBrush)
	canvas.DrawRoundedRect(s.barRect, 3, 3, 3, 3, s.barPen, s.barBrush)
}

func (s *ScrollBarImpl) RailBrush() Brush {
	return s.railBrush
}

func (s *ScrollBarImpl) SetRailBrush(brush Brush) {
	if s.railBrush != brush {
		s.railBrush = brush
		s.Redraw()
	}
}

func (s *ScrollBarImpl) BarBrush() Brush {
	return s.barBrush
}

func (s *ScrollBarImpl) SetBarBrush(brush Brush) {
	if s.barBrush != brush {
		s.barBrush = brush
		s.Redraw()
	}
}

func (s *ScrollBarImpl) RailPen() Pen {
	return s.railPen
}

func (s *ScrollBarImpl) SetRailPen(pen Pen) {
	if s.railPen != pen {
		s.railPen = pen
		s.Redraw()
	}
}

func (s *ScrollBarImpl) BarPen() Pen {
	return s.barPen
}

func (s *ScrollBarImpl) SetBarPen(pen Pen) {
	if s.barPen != pen {
		s.barPen = pen
		s.Redraw()
	}
}

func (s *ScrollBarImpl) ScrollPosition() (int, int) {
	return s.scrollPositionFrom, s.scrollPositionTo
}

func (s *ScrollBarImpl) SetScrollPosition(from, to int) {
	if s.scrollPositionFrom != from || s.scrollPositionTo != to {
		s.scrollPositionFrom, s.scrollPositionTo = from, to
		s.updateBarRect()
		s.Redraw()
		s.onScroll.Fire(from, to)
	}
}

func (s *ScrollBarImpl) ScrollLimit() int {
	return s.scrollLimit
}

func (s *ScrollBarImpl) SetScrollLimit(limit int) {
	if s.scrollLimit != limit {
		s.scrollLimit = limit
		s.updateBarRect()
		s.Redraw()
	}
}

func (s *ScrollBarImpl) AutoHide() bool {
	return s.autoHide
}

func (s *ScrollBarImpl) SetAutoHide(autoHide bool) {
	if s.autoHide != autoHide {
		s.autoHide = autoHide
		s.Redraw()
	}
}

func (s *ScrollBarImpl) IsVisible() bool {
	if s.autoHide && s.scrollPositionFrom == 0 && s.scrollPositionTo == s.scrollLimit {
		return false
	}
	return s.ControlBase.IsVisible()
}

func (s *ScrollBarImpl) Orientation() Orientation {
	return s.orientation
}

func (s *ScrollBarImpl) SetOrientation(orientation Orientation) {
	if s.orientation != orientation {
		s.orientation = orientation
		s.Redraw()
	}
}

// InputEventHandlerPart overrides
func (s *ScrollBarImpl) Click(event MouseEvent) bool {
	if !s.barRect.Contains(event.Point) {
		p := s.positionAt(event.Point)
		from, to := s.scrollPositionFrom, s.scrollPositionTo
		switch {
		case p < from:
			width := to - from
			from = math.Max(from-width, 0)
			s.SetScrollPosition(from, from+width)
		case p > to:
			width := to - from
			to = math.Min(to+width, s.scrollLimit)
			s.SetScrollPosition(to-width, to)
		}
	}
	return true
}

func (s *ScrollBarImpl) MouseDown(event MouseEvent) {
	if s.barRect.Contains(event.Point) {
		initialOffset := event.Point.Sub(s.barRect.Min)
		var mms, mus EventSubscription
		mms = event.Window.OnMouseMove(func(we MouseEvent) {
			p := WindowToChild(we.WindowPoint, s.parent)
			s.SetScrollPosition(s.rangeAt(p.Sub(initialOffset)))
		})
		mus = event.Window.OnMouseUp(func(we MouseEvent) {
			mms.Forget()
			mus.Forget()
		})
	}
	s.InputEventHandlerPart.MouseDown(event)
}
