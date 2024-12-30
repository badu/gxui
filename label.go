// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
	"strings"
)

type Label interface {
	Control
	Text() string
	SetText(text string)
	Font() Font
	SetFont(font Font)
	Color() Color
	SetColor(color Color)
	Multiline() bool
	SetMultiline(bool)
	SetHorizontalAlignment(HorizontalAlignment)
	HorizontalAlignment() HorizontalAlignment
	SetVerticalAlignment(VerticalAlignment)
	VerticalAlignment() VerticalAlignment
}

type LabelImpl struct {
	ControlBase
	outer               ControlBaseOuter
	font                Font
	color               Color
	horizontalAlignment HorizontalAlignment
	verticalAlignment   VerticalAlignment
	multiline           bool
	text                string
}

func (l *LabelImpl) Init(outer ControlBaseOuter, theme Theme, font Font, color Color) {
	if font == nil {
		panic("Cannot create a label with a nil font")
	}
	l.ControlBase.Init(outer, theme)
	l.outer = outer
	l.font = font
	l.color = color
	l.horizontalAlignment = AlignLeft
	l.verticalAlignment = AlignMiddle
}

func (l *LabelImpl) Text() string {
	return l.text
}

func (l *LabelImpl) SetText(text string) {
	if l.text != text {
		l.text = text
		l.outer.Relayout()
	}
}

func (l *LabelImpl) Font() Font {
	return l.font
}

func (l *LabelImpl) SetFont(font Font) {
	if l.font != font {
		l.font = font
		l.Relayout()
	}
}

func (l *LabelImpl) Color() Color {
	return l.color
}

func (l *LabelImpl) SetColor(color Color) {
	if l.color != color {
		l.color = color
		l.outer.Redraw()
	}
}

func (l *LabelImpl) Multiline() bool {
	return l.multiline
}

func (l *LabelImpl) SetMultiline(multiline bool) {
	if l.multiline != multiline {
		l.multiline = multiline
		l.outer.Relayout()
	}
}

func (l *LabelImpl) DesiredSize(min, max math.Size) math.Size {
	text := l.text
	if !l.multiline {
		text = strings.Replace(text, "\n", " ", -1)
	}
	size := l.font.Measure(&TextBlock{Runes: []rune(text)})
	return size.Clamp(min, max)
}

func (l *LabelImpl) SetHorizontalAlignment(horizontalAlignment HorizontalAlignment) {
	if l.horizontalAlignment != horizontalAlignment {
		l.horizontalAlignment = horizontalAlignment
		l.Redraw()
	}
}

func (l *LabelImpl) HorizontalAlignment() HorizontalAlignment {
	return l.horizontalAlignment
}

func (l *LabelImpl) SetVerticalAlignment(verticalAlignment VerticalAlignment) {
	if l.verticalAlignment != verticalAlignment {
		l.verticalAlignment = verticalAlignment
		l.Redraw()
	}
}

func (l *LabelImpl) VerticalAlignment() VerticalAlignment {
	return l.verticalAlignment
}

// parts.DrawPaintPart overrides
func (l *LabelImpl) Paint(canvas Canvas) {
	rect := l.outer.Size().Rect()
	text := l.text
	if !l.multiline {
		text = strings.Replace(text, "\n", " ", -1)
	}

	runes := []rune(text)
	offsets := l.font.Layout(&TextBlock{
		Runes:     runes,
		AlignRect: rect,
		H:         l.horizontalAlignment,
		V:         l.verticalAlignment,
	})
	canvas.DrawRunes(l.font, runes, offsets, l.color)
}
