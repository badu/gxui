// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
	"strings"
)

type Label struct {
	ControlBase
	parent              ControlBaseParent
	font                Font
	color               Color
	horizontalAlignment HAlign
	verticalAlignment   VAlign
	text                string
	multiline           bool
}

func (l *Label) Init(parent ControlBaseParent, driver Driver, styles *StyleDefs) {
	l.ControlBase.Init(parent, driver)
	l.parent = parent
	l.font = styles.DefaultFont
	l.color = styles.LabelStyle.FontColor
	l.horizontalAlignment = styles.LabelStyle.HAlign
	l.verticalAlignment = styles.LabelStyle.VAlign
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) SetText(text string) {
	if l.text == text {
		return
	}

	l.text = text
	l.parent.ReLayout()
}

func (l *Label) Font() Font {
	return l.font
}

func (l *Label) SetFont(font Font) {
	if l.font == font {
		return
	}

	l.font = font
	l.ReLayout()
}

func (l *Label) Color() Color {
	return l.color
}

func (l *Label) SetColor(color Color) {
	if l.color == color {
		return
	}

	l.color = color
	l.parent.Redraw()
}

func (l *Label) Multiline() bool {
	return l.multiline
}

func (l *Label) SetMultiline(multiline bool) {
	if l.multiline == multiline {
		return
	}

	l.multiline = multiline
	l.parent.ReLayout()
}

func (l *Label) DesiredSize(min, max math.Size) math.Size {
	text := l.text
	if !l.multiline {
		text = strings.Replace(text, "\n", " ", -1)
	}
	size := l.font.Measure(&TextBlock{Runes: []rune(text)})
	return size.Clamp(min, max)
}

func (l *Label) SetHorizontalAlignment(horizontalAlignment HAlign) {
	if l.horizontalAlignment == horizontalAlignment {
		return
	}

	l.horizontalAlignment = horizontalAlignment
	l.Redraw()
}

func (l *Label) HorizontalAlignment() HAlign {
	return l.horizontalAlignment
}

func (l *Label) SetVerticalAlignment(verticalAlignment VAlign) {
	if l.verticalAlignment != verticalAlignment {
		return
	}
	l.verticalAlignment = verticalAlignment
	l.Redraw()
}

func (l *Label) VerticalAlignment() VAlign {
	return l.verticalAlignment
}

// parts.DrawPaintPart overrides
func (l *Label) Paint(canvas Canvas) {
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
