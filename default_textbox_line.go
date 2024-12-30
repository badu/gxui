// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/interval"
	"github.com/badu/gxui/math"
)

// DefaultTextBoxLine
type DefaultTextBoxLine struct {
	ControlBase
	parent     DefaultTextBoxLineParent
	textbox    *TextBoxImpl
	lineIndex  int
	caretWidth int
}

func (t *DefaultTextBoxLine) Init(parent DefaultTextBoxLineParent, app App, textbox *TextBoxImpl, lineIndex int) {
	t.ControlBase.Init(parent, app)
	t.parent = parent
	t.textbox = textbox
	t.lineIndex = lineIndex
	t.SetCaretWidth(2)
	t.OnAttach(func() {
		ev := t.textbox.OnRedrawLines(t.Redraw)
		t.OnDetach(ev.Forget)
	})
}

func (t *DefaultTextBoxLine) SetCaretWidth(width int) {
	if t.caretWidth != width {
		t.caretWidth = width
	}
}

func (t *DefaultTextBoxLine) DesiredSize(min, max math.Size) math.Size {
	return max
}

func (t *DefaultTextBoxLine) Paint(canvas Canvas) {
	if t.textbox.HasFocus() {
		t.parent.PaintSelections(canvas)
	}

	t.parent.PaintText(canvas)

	if t.textbox.HasFocus() {
		t.parent.PaintCarets(canvas)
	}
}

func (t *DefaultTextBoxLine) MeasureRunes(start, end int) math.Size {
	controller := t.textbox.controller
	return t.textbox.font.Measure(&TextBlock{
		Runes: controller.TextRunes()[start:end],
	})
}

func (t *DefaultTextBoxLine) PaintText(canvas Canvas) {
	runes := []rune(t.textbox.controller.Line(t.lineIndex))
	textFont := t.textbox.font
	offsets := textFont.Layout(&TextBlock{
		Runes:     runes,
		AlignRect: t.Size().Rect().OffsetX(t.caretWidth),
		H:         AlignLeft,
		V:         AlignBottom,
	})
	canvas.DrawRunes(textFont, runes, offsets, t.textbox.textColor)
}

func (t *DefaultTextBoxLine) PaintCarets(canvas Canvas) {
	controller := t.textbox.controller
	for i, cnt := 0, controller.SelectionCount(); i < cnt; i++ {
		caretEnd := controller.Caret(i)
		lineIndex := controller.LineIndex(caretEnd)
		if lineIndex == t.lineIndex {
			start := controller.LineStart(lineIndex)
			measuredRunes := t.parent.MeasureRunes(start, caretEnd)
			top := math.Point{X: t.caretWidth + measuredRunes.W, Y: 0}
			bottom := top.Add(math.Point{X: 0, Y: t.Size().H})
			t.parent.PaintCaret(canvas, top, bottom)
		}
	}
}

func (t *DefaultTextBoxLine) PaintSelections(canvas Canvas) {
	controller := t.textbox.controller

	lineStart, lineEnd := controller.LineStart(t.lineIndex), controller.LineEnd(t.lineIndex)

	selections := controller.Selections()
	if t.textbox.selectionDragging {
		interval.Replace(&selections, t.textbox.selectionDrag)
	}

	interval.Visit(
		&selections,
		CreateTextSelection(lineStart, lineEnd, false),
		func(s, e uint64, _ int) {
			if s < e {
				x := t.parent.MeasureRunes(lineStart, int(s)).W
				m := t.parent.MeasureRunes(int(s), int(e))
				top := math.Point{X: t.caretWidth + x, Y: 0}
				bottom := top.Add(m.Point())
				t.parent.PaintSelection(canvas, top, bottom)
			}
		},
	)
}

func (t *DefaultTextBoxLine) PaintCaret(canvas Canvas, top, bottom math.Point) {
	rect := math.Rect{Min: top, Max: bottom}.ExpandI(t.caretWidth / 2)
	canvas.DrawRoundedRect(rect, 1, 1, 1, 1, CreatePen(0.5, Gray70), WhiteBrush)
}

func (t *DefaultTextBoxLine) PaintSelection(canvas Canvas, top, bottom math.Point) {
	rect := math.Rect{Min: top, Max: bottom}.ExpandI(t.caretWidth / 2)
	canvas.DrawRoundedRect(rect, 1, 1, 1, 1, TransparentPen, Brush{Color: Gray40})
}

// TextBoxLine compliance
func (t *DefaultTextBoxLine) RuneIndexAt(point math.Point) int {
	font := t.textbox.font
	controller := t.textbox.controller

	x := point.X
	line := controller.Line(t.lineIndex)
	i := 0
	for ; i < len(line) && x > font.Measure(&TextBlock{Runes: []rune(line[:i+1])}).W; i++ {
	}

	return controller.LineStart(t.lineIndex) + i
}

func (t *DefaultTextBoxLine) PositionAt(runeIndex int) math.Point {
	font := t.textbox.font
	controller := t.textbox.controller

	x := runeIndex - controller.LineStart(t.lineIndex)
	line := controller.Line(t.lineIndex)
	return font.Measure(&TextBlock{Runes: []rune(line[:x])}).Point()
}
