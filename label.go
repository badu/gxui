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
	SetHorizontalAlignment(HAlign)
	HorizontalAlignment() HAlign
	SetVerticalAlignment(VAlign)
	VerticalAlignment() VAlign
}

type LabelImpl struct {
	ControlBase
	parent              ControlBaseParent
	font                Font
	color               Color
	horizontalAlignment HAlign
	verticalAlignment   VAlign
	text                string
	multiline           bool
}

func (l *LabelImpl) Init(parent ControlBaseParent, driver Driver, styles *StyleDefs) {
	l.ControlBase.Init(parent, driver)
	l.parent = parent
	l.font = styles.DefaultFont
	l.color = styles.LabelStyle.FontColor
	l.horizontalAlignment = styles.LabelStyle.HAlign
	l.verticalAlignment = styles.LabelStyle.VAlign
}

func (l *LabelImpl) Text() string {
	return l.text
}

func (l *LabelImpl) SetText(text string) {
	if l.text == text {
		return
	}

	l.text = text
	l.parent.ReLayout()
}

func (l *LabelImpl) Font() Font {
	return l.font
}

func (l *LabelImpl) SetFont(font Font) {
	if l.font == font {
		return
	}

	l.font = font
	l.ReLayout()
}

func (l *LabelImpl) Color() Color {
	return l.color
}

func (l *LabelImpl) SetColor(color Color) {
	if l.color == color {
		return
	}

	l.color = color
	l.parent.Redraw()
}

func (l *LabelImpl) Multiline() bool {
	return l.multiline
}

func (l *LabelImpl) SetMultiline(multiline bool) {
	if l.multiline == multiline {
		return
	}

	l.multiline = multiline
	l.parent.ReLayout()
}

func (l *LabelImpl) DesiredSize(min, max math.Size) math.Size {
	text := l.text
	if !l.multiline {
		text = strings.Replace(text, "\n", " ", -1)
	}
	size := l.font.Measure(&TextBlock{Runes: []rune(text)})
	return size.Clamp(min, max)
}

func (l *LabelImpl) SetHorizontalAlignment(horizontalAlignment HAlign) {
	if l.horizontalAlignment == horizontalAlignment {
		return
	}

	l.horizontalAlignment = horizontalAlignment
	l.Redraw()
}

func (l *LabelImpl) HorizontalAlignment() HAlign {
	return l.horizontalAlignment
}

func (l *LabelImpl) SetVerticalAlignment(verticalAlignment VAlign) {
	if l.verticalAlignment != verticalAlignment {
		return
	}
	l.verticalAlignment = verticalAlignment
	l.Redraw()
}

func (l *LabelImpl) VerticalAlignment() VAlign {
	return l.verticalAlignment
}

// parts.DrawPaintPart overrides
func (l *LabelImpl) Paint(canvas Canvas) {
	rect := l.parent.Size().Rect()
	text := l.text
	if !l.multiline {
		text = strings.Replace(text, "\n", " ", -1)
	}

	runes := []rune(text)
	offsets := l.font.Layout(
		&TextBlock{Runes: runes, AlignRect: rect, H: l.horizontalAlignment, V: l.verticalAlignment},
	)
	canvas.DrawRunes(l.font, runes, offsets, l.color)
}
