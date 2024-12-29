// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/mixins/base"
)

type ScrollBarOuter interface {
	base.ControlOuter
}

type ScrollBar struct {
	base.Control
	outer               ScrollBarOuter
	orientation         gxui.Orientation
	thickness           int
	minBarLength        int
	scrollPositionFrom  int
	scrollPositionTo    int
	scrollLimit         int
	railBrush, barBrush gxui.Brush
	railPen, barPen     gxui.Pen
	barRect             math.Rect
	onScroll            gxui.Event
	autoHide            bool
}

func (s *ScrollBar) positionAt(p math.Point) int {
	orientation := s.orientation
	fraction := float32(orientation.Major(p.XY())) / float32(orientation.Major(s.Size().WH()))
	limit := s.ScrollLimit()
	return int(float32(limit) * fraction)
}

func (s *ScrollBar) rangeAt(point math.Point) (int, int) {
	width := s.scrollPositionTo - s.scrollPositionFrom
	from := math.Clamp(s.positionAt(point), 0, s.scrollLimit-width)
	to := from + width
	return from, to
}

func (s *ScrollBar) updateBarRect() {
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

func (s *ScrollBar) Init(outer ScrollBarOuter, theme gxui.Theme) {
	s.Control.Init(outer, theme)

	s.outer = outer
	s.thickness = 10
	s.minBarLength = 10
	s.scrollPositionFrom = 0
	s.scrollPositionTo = 100
	s.scrollLimit = 100
	s.onScroll = gxui.CreateEvent(s.SetScrollPosition)
}

func (s *ScrollBar) OnScroll(callback func(from, to int)) gxui.EventSubscription {
	return s.onScroll.Listen(callback)
}

func (s *ScrollBar) ScrollFraction() (float32, float32) {
	from := float32(s.scrollPositionFrom) / float32(s.scrollLimit)
	to := float32(s.scrollPositionTo) / float32(s.scrollLimit)
	return from, to
}

func (s *ScrollBar) DesiredSize(min, max math.Size) math.Size {
	if s.orientation.Horizontal() {
		return math.Size{W: max.W, H: s.thickness}.Clamp(min, max)
	} else {
		return math.Size{W: s.thickness, H: max.H}.Clamp(min, max)
	}
}

func (s *ScrollBar) Paint(canvas gxui.Canvas) {
	canvas.DrawRoundedRect(s.outer.Size().Rect(), 3, 3, 3, 3, s.railPen, s.railBrush)
	canvas.DrawRoundedRect(s.barRect, 3, 3, 3, 3, s.barPen, s.barBrush)
}

func (s *ScrollBar) RailBrush() gxui.Brush {
	return s.railBrush
}

func (s *ScrollBar) SetRailBrush(brush gxui.Brush) {
	if s.railBrush != brush {
		s.railBrush = brush
		s.Redraw()
	}
}

func (s *ScrollBar) BarBrush() gxui.Brush {
	return s.barBrush
}

func (s *ScrollBar) SetBarBrush(brush gxui.Brush) {
	if s.barBrush != brush {
		s.barBrush = brush
		s.Redraw()
	}
}

func (s *ScrollBar) RailPen() gxui.Pen {
	return s.railPen
}

func (s *ScrollBar) SetRailPen(pen gxui.Pen) {
	if s.railPen != pen {
		s.railPen = pen
		s.Redraw()
	}
}

func (s *ScrollBar) BarPen() gxui.Pen {
	return s.barPen
}

func (s *ScrollBar) SetBarPen(pen gxui.Pen) {
	if s.barPen != pen {
		s.barPen = pen
		s.Redraw()
	}
}

func (s *ScrollBar) ScrollPosition() (int, int) {
	return s.scrollPositionFrom, s.scrollPositionTo
}

func (s *ScrollBar) SetScrollPosition(from, to int) {
	if s.scrollPositionFrom != from || s.scrollPositionTo != to {
		s.scrollPositionFrom, s.scrollPositionTo = from, to
		s.updateBarRect()
		s.Redraw()
		s.onScroll.Fire(from, to)
	}
}

func (s *ScrollBar) ScrollLimit() int {
	return s.scrollLimit
}

func (s *ScrollBar) SetScrollLimit(limit int) {
	if s.scrollLimit != limit {
		s.scrollLimit = limit
		s.updateBarRect()
		s.Redraw()
	}
}

func (s *ScrollBar) AutoHide() bool {
	return s.autoHide
}

func (s *ScrollBar) SetAutoHide(autoHide bool) {
	if s.autoHide != autoHide {
		s.autoHide = autoHide
		s.Redraw()
	}
}

func (s *ScrollBar) IsVisible() bool {
	if s.autoHide && s.scrollPositionFrom == 0 && s.scrollPositionTo == s.scrollLimit {
		return false
	}
	return s.Control.IsVisible()
}

func (s *ScrollBar) Orientation() gxui.Orientation {
	return s.orientation
}

func (s *ScrollBar) SetOrientation(orientation gxui.Orientation) {
	if s.orientation != orientation {
		s.orientation = orientation
		s.Redraw()
	}
}

// InputEventHandler overrides
func (s *ScrollBar) Click(event gxui.MouseEvent) bool {
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

func (s *ScrollBar) MouseDown(event gxui.MouseEvent) {
	if s.barRect.Contains(event.Point) {
		initialOffset := event.Point.Sub(s.barRect.Min)
		var mms, mus gxui.EventSubscription
		mms = event.Window.OnMouseMove(func(we gxui.MouseEvent) {
			p := gxui.WindowToChild(we.WindowPoint, s.outer)
			s.SetScrollPosition(s.rangeAt(p.Sub(initialOffset)))
		})
		mus = event.Window.OnMouseUp(func(we gxui.MouseEvent) {
			mms.Unlisten()
			mus.Unlisten()
		})
	}
	s.InputEventHandler.MouseDown(event)
}
