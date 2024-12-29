// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
)

type TextBox interface {
	Focusable
	OnSelectionChanged(func()) EventSubscription
	OnTextChanged(func([]TextBoxEdit)) EventSubscription
	Padding() math.Spacing
	SetPadding(math.Spacing)
	Runes() []rune
	Text() string
	SetText(string)
	Font() Font
	SetFont(Font)
	Multiline() bool
	SetMultiline(bool)
	DesiredWidth() int
	SetDesiredWidth(desiredWidth int)
	TextColor() Color
	SetTextColor(Color)
	Select(TextSelectionList)
	SelectAll()
	Carets() []int
	RuneIndexAt(p math.Point) (idx int, found bool)
	TextAt(s, e int) string
	WordAt(runeIndex int) string
	ScrollToLine(int)
	ScrollToRune(int)
	LineIndex(runeIndex int) int
	LineStart(line int) int
	LineEnd(line int) int
}

type TextBoxLine interface {
	Control
	RuneIndexAt(math.Point) int
	PositionAt(int) math.Point
}

type TextBoxOuter interface {
	ListOuter
	CreateLine(theme Theme, index int) (line TextBoxLine, container Control)
}

type DefaultTextBoxLineOuter interface {
	ControlBaseOuter
	MeasureRunes(s, e int) math.Size
	PaintText(c Canvas)
	PaintCarets(c Canvas)
	PaintCaret(c Canvas, top, bottom math.Point)
	PaintSelections(c Canvas)
	PaintSelection(c Canvas, top, bottom math.Point)
}
