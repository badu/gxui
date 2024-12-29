// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/interval"
	"github.com/badu/gxui/math"
)

// DefaultTextBoxLine
type DefaultTextBoxLine struct {
	ControlBase
	outer      gxui.DefaultTextBoxLineOuter
	textbox    *TextBox
	lineIndex  int
	caretWidth int
}

func (t *DefaultTextBoxLine) Init(outer gxui.DefaultTextBoxLineOuter, theme gxui.Theme, textbox *TextBox, lineIndex int) {
	t.ControlBase.Init(outer, theme)
	t.outer = outer
	t.textbox = textbox
	t.lineIndex = lineIndex
	t.SetCaretWidth(2)
	t.OnAttach(func() {
		ev := t.textbox.OnRedrawLines(t.Redraw)
		t.OnDetach(ev.Unlisten)
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

func (t *DefaultTextBoxLine) Paint(canvas gxui.Canvas) {
	if t.textbox.HasFocus() {
		t.outer.PaintSelections(canvas)
	}

	t.outer.PaintText(canvas)

	if t.textbox.HasFocus() {
		t.outer.PaintCarets(canvas)
	}
}

func (t *DefaultTextBoxLine) MeasureRunes(start, end int) math.Size {
	controller := t.textbox.controller
	return t.textbox.font.Measure(&gxui.TextBlock{
		Runes: controller.TextRunes()[start:end],
	})
}

func (t *DefaultTextBoxLine) PaintText(canvas gxui.Canvas) {
	runes := []rune(t.textbox.controller.Line(t.lineIndex))
	textFont := t.textbox.font
	offsets := textFont.Layout(&gxui.TextBlock{
		Runes:     runes,
		AlignRect: t.Size().Rect().OffsetX(t.caretWidth),
		H:         gxui.AlignLeft,
		V:         gxui.AlignBottom,
	})
	canvas.DrawRunes(textFont, runes, offsets, t.textbox.textColor)
}

func (t *DefaultTextBoxLine) PaintCarets(canvas gxui.Canvas) {
	controller := t.textbox.controller
	for i, cnt := 0, controller.SelectionCount(); i < cnt; i++ {
		caretEnd := controller.Caret(i)
		lineIndex := controller.LineIndex(caretEnd)
		if lineIndex == t.lineIndex {
			start := controller.LineStart(lineIndex)
			measuredRunes := t.outer.MeasureRunes(start, caretEnd)
			top := math.Point{X: t.caretWidth + measuredRunes.W, Y: 0}
			bottom := top.Add(math.Point{X: 0, Y: t.Size().H})
			t.outer.PaintCaret(canvas, top, bottom)
		}
	}
}

func (t *DefaultTextBoxLine) PaintSelections(canvas gxui.Canvas) {
	controller := t.textbox.controller

	lineStart, lineEnd := controller.LineStart(t.lineIndex), controller.LineEnd(t.lineIndex)

	selections := controller.Selections()
	if t.textbox.selectionDragging {
		interval.Replace(&selections, t.textbox.selectionDrag)
	}

	interval.Visit(
		&selections,
		gxui.CreateTextSelection(lineStart, lineEnd, false),
		func(s, e uint64, _ int) {
			if s < e {
				x := t.outer.MeasureRunes(lineStart, int(s)).W
				m := t.outer.MeasureRunes(int(s), int(e))
				top := math.Point{X: t.caretWidth + x, Y: 0}
				bottom := top.Add(m.Point())
				t.outer.PaintSelection(canvas, top, bottom)
			}
		},
	)
}

func (t *DefaultTextBoxLine) PaintCaret(canvas gxui.Canvas, top, bottom math.Point) {
	rect := math.Rect{Min: top, Max: bottom}.ExpandI(t.caretWidth / 2)
	canvas.DrawRoundedRect(rect, 1, 1, 1, 1, gxui.CreatePen(0.5, gxui.Gray70), gxui.WhiteBrush)
}

func (t *DefaultTextBoxLine) PaintSelection(canvas gxui.Canvas, top, bottom math.Point) {
	rect := math.Rect{Min: top, Max: bottom}.ExpandI(t.caretWidth / 2)
	canvas.DrawRoundedRect(rect, 1, 1, 1, 1, gxui.TransparentPen, gxui.Brush{Color: gxui.Gray40})
}

// TextBoxLine compliance
func (t *DefaultTextBoxLine) RuneIndexAt(point math.Point) int {
	font := t.textbox.font
	controller := t.textbox.controller

	x := point.X
	line := controller.Line(t.lineIndex)
	i := 0
	for ; i < len(line) && x > font.Measure(&gxui.TextBlock{Runes: []rune(line[:i+1])}).W; i++ {
	}

	return controller.LineStart(t.lineIndex) + i
}

func (t *DefaultTextBoxLine) PositionAt(runeIndex int) math.Point {
	font := t.textbox.font
	controller := t.textbox.controller

	x := runeIndex - controller.LineStart(t.lineIndex)
	line := controller.Line(t.lineIndex)
	return font.Measure(&gxui.TextBlock{Runes: []rune(line[:x])}).Point()
}
