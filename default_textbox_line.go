// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/interval"
	"github.com/badu/gxui/pkg/math"
)

type DefaultTextBoxLineParent interface {
	// Control interface
	Size() math.Size
	SetSize(newSize math.Size)
	Draw() Canvas
	Parent() Parent
	SetParent(newParent Parent)
	Attached() bool
	Attach()
	Detach()
	DesiredSize(min, max math.Size) math.Size
	Margin() math.Spacing
	SetMargin(math.Spacing)
	IsVisible() bool
	SetVisible(isVisible bool)
	ContainsPoint(point math.Point) bool
	IsMouseOver() bool
	IsMouseDown(button MouseButton) bool
	Click(event MouseEvent) (consume bool)
	DoubleClick(event MouseEvent) (consume bool)
	KeyPress(event KeyboardEvent) (consume bool)
	KeyStroke(event KeyStrokeEvent) (consume bool)
	MouseScroll(event MouseEvent) (consume bool)
	MouseMove(event MouseEvent)
	MouseEnter(event MouseEvent)
	MouseExit(event MouseEvent)
	MouseDown(event MouseEvent)
	MouseUp(event MouseEvent)
	KeyDown(event KeyboardEvent)
	KeyUp(event KeyboardEvent)
	KeyRepeat(event KeyboardEvent)
	OnAttach(callback func()) EventSubscription
	OnDetach(callback func()) EventSubscription
	OnKeyPress(callback func(KeyboardEvent)) EventSubscription
	OnKeyStroke(callback func(KeyStrokeEvent)) EventSubscription
	OnClick(callback func(MouseEvent)) EventSubscription
	OnDoubleClick(callback func(MouseEvent)) EventSubscription
	OnMouseMove(callback func(MouseEvent)) EventSubscription
	OnMouseEnter(callback func(MouseEvent)) EventSubscription
	OnMouseExit(callback func(MouseEvent)) EventSubscription
	OnMouseDown(callback func(MouseEvent)) EventSubscription
	OnMouseUp(callback func(MouseEvent)) EventSubscription
	OnMouseScroll(callback func(MouseEvent)) EventSubscription
	OnKeyDown(callback func(KeyboardEvent)) EventSubscription
	OnKeyUp(callback func(KeyboardEvent)) EventSubscription
	OnKeyRepeat(callback func(KeyboardEvent)) EventSubscription
	// ControlBaseParent interface
	Paint(canvas Canvas) // was outer.Painter
	Redraw()             // was outer.Redrawer
	ReLayout()           // was outer.Relayouter
	MeasureRunes(s, e int) math.Size
	PaintText(c Canvas)
	PaintCarets(c Canvas)
	PaintCaret(c Canvas, top, bottom math.Point)
	PaintSelections(c Canvas)
	PaintSelection(c Canvas, top, bottom math.Point)
}

// DefaultTextBoxLine
type DefaultTextBoxLine struct {
	InputEventHandlerPart
	ParentablePart
	DrawPaintPart
	AttachablePart
	VisiblePart
	LayoutablePart
	parent     DefaultTextBoxLineParent
	textbox    *TextBox
	lineIndex  int
	caretWidth int
	offset     int
}

func (t *DefaultTextBoxLine) Init(parent DefaultTextBoxLineParent, textbox *TextBox, lineIndex int) {
	t.DrawPaintPart.Init(parent, textbox.driver)
	t.LayoutablePart.Init(parent)
	t.InputEventHandlerPart.Init()
	t.VisiblePart.Init(parent)

	t.parent = parent
	t.textbox = textbox
	t.lineIndex = lineIndex
	t.SetCaretWidth(2)
	t.OnAttach(
		func() {
			ev := t.textbox.OnRedrawLines(t.Redraw)
			t.OnDetach(ev.Forget)
		},
	)
}

func (t *DefaultTextBoxLine) SetCaretWidth(width int) {
	if t.caretWidth == width {
		return
	}

	t.caretWidth = width
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
	return t.textbox.font.Measure(
		&TextBlock{Runes: controller.TextRunes()[start:end]},
	)
}

func (t *DefaultTextBoxLine) PaintText(canvas Canvas) {
	runes := []rune(t.textbox.controller.Line(t.lineIndex))
	textFont := t.textbox.font
	offsets := textFont.Layout(
		&TextBlock{
			Runes:     runes,
			AlignRect: t.Size().Rect().OffsetX(t.caretWidth),
			H:         AlignLeft,
			V:         AlignBottom,
		},
	)
	canvas.DrawRunes(textFont, runes, offsets, t.textbox.textColor)
}

func (t *DefaultTextBoxLine) PaintCarets(canvas Canvas) {
	controller := t.textbox.controller
	for caret, count := 0, controller.SelectionCount(); caret < count; caret++ {
		caretEnd := controller.Caret(caret)
		lineIndex := controller.LineIndex(caretEnd)

		if lineIndex != t.lineIndex {
			continue
		}

		start := controller.LineStart(lineIndex)
		measuredRunes := t.parent.MeasureRunes(start, caretEnd)
		top := math.Point{X: t.caretWidth + measuredRunes.Width, Y: 0}
		bottom := top.Add(math.Point{X: 0, Y: t.Size().Height})
		t.parent.PaintCaret(canvas, top, bottom)
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
				x := t.parent.MeasureRunes(lineStart, int(s)).Width
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

	// TODO : @Badu - why?
	for ; i < len(line) && x > font.Measure(&TextBlock{Runes: []rune(line[:i+1])}).Width; i++ {
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

func (t *DefaultTextBoxLine) SetOffset(offset int) {
	t.offset = offset
	t.Redraw()
}

func (t *DefaultTextBoxLine) ContainsPoint(point math.Point) bool {
	return t.IsVisible() && t.Size().Rect().Contains(point)
}
