// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/math"
	"strings"
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

type TextBoxParent interface {
	ListParent
	CreateLine(driver Driver, styles *StyleDefs, index int) (line TextBoxLine, container Control)
}

type DefaultTextBoxLineParent interface {
	ControlBaseParent
	MeasureRunes(s, e int) math.Size
	PaintText(c Canvas)
	PaintCarets(c Canvas)
	PaintCaret(c Canvas, top, bottom math.Point)
	PaintSelections(c Canvas)
	PaintSelection(c Canvas, top, bottom math.Point)
}

type TextBoxImpl struct {
	ListImpl
	AdapterBase
	FocusablePart
	parent            TextBoxParent
	driver            Driver
	font              Font
	textColor         Color
	onRedrawLines     Event
	multiline         bool
	controller        *TextBoxController
	adapter           *TextBoxAdapter
	selectionDragging bool
	selectionDrag     TextSelection
	desiredWidth      int
}

func (t *TextBoxImpl) lineMouseDown(line TextBoxLine, event MouseEvent) {
	if event.Button == MouseButtonLeft {
		p := line.RuneIndexAt(event.Point)
		t.selectionDragging = true
		t.selectionDrag = CreateTextSelection(p, p, false)
		if !event.Modifier.Control() {
			t.controller.SetCaret(p)
		}
	}
}

func (t *TextBoxImpl) lineMouseUp(line TextBoxLine, event MouseEvent) {
	if event.Button == MouseButtonLeft {
		t.selectionDragging = false
		if !event.Modifier.Control() {
			t.controller.SetSelection(t.selectionDrag)
		} else {
			t.controller.AddSelection(t.selectionDrag)
		}
	}
}

func (t *TextBoxImpl) Init(parent TextBoxParent, driver Driver, styles *StyleDefs, font Font) {
	t.ListImpl.Init(parent, driver, styles)
	t.FocusablePart.Init()
	t.parent = parent
	t.driver = driver
	if font == nil {
		t.font = styles.DefaultFont
	} else {
		t.font = font
	}
	t.onRedrawLines = CreateEvent(func() {})
	t.controller = CreateTextBoxController()
	t.adapter = &TextBoxAdapter{TextBox: t}
	t.desiredWidth = 100
	t.SetScrollBarEnabled(false) // Defaults to single line
	t.OnGainedFocus(func() { t.onRedrawLines.Fire() })
	t.OnLostFocus(func() { t.onRedrawLines.Fire() })
	t.controller.OnTextChanged(func([]TextBoxEdit) {
		t.onRedrawLines.Fire()
		t.ListImpl.DataChanged(false)
	})
	t.controller.OnSelectionChanged(func() {
		t.onRedrawLines.Fire()
	})

	t.ListImpl.SetAdapter(t.adapter)
}

func (t *TextBoxImpl) textRect() math.Rect {
	return t.parent.Size().Rect().Contract(t.Padding())
}

func (t *TextBoxImpl) pageLines() int {
	return (t.parent.Size().H - t.parent.Padding().H()) / t.MajorAxisItemSize()
}

func (t *TextBoxImpl) OnRedrawLines(callback func()) EventSubscription {
	return t.onRedrawLines.Listen(callback)
}

func (t *TextBoxImpl) OnSelectionChanged(callback func()) EventSubscription {
	return t.controller.OnSelectionChanged(callback)
}

func (t *TextBoxImpl) OnTextChanged(callback func(lines []TextBoxEdit)) EventSubscription {
	return t.controller.OnTextChanged(callback)
}

func (t *TextBoxImpl) Runes() []rune {
	return t.controller.TextRunes()
}

func (t *TextBoxImpl) Text() string {
	return t.controller.Text()
}

func (t *TextBoxImpl) SetText(text string) {
	t.controller.SetText(text)
	t.parent.Relayout()
}

func (t *TextBoxImpl) TextColor() Color {
	return t.textColor
}

func (t *TextBoxImpl) SetTextColor(color Color) {
	t.textColor = color
	t.Relayout()
}

func (t *TextBoxImpl) Font() Font {
	return t.font
}

func (t *TextBoxImpl) SetFont(font Font) {
	if t.font != font {
		t.font = font
		t.Relayout()
	}
}

func (t *TextBoxImpl) Multiline() bool {
	return t.multiline
}

func (t *TextBoxImpl) SetMultiline(multiline bool) {
	if t.multiline != multiline {
		t.multiline = multiline
		t.SetScrollBarEnabled(multiline)
		t.parent.Relayout()
	}
}

func (t *TextBoxImpl) DesiredWidth() int {
	return t.desiredWidth
}

func (t *TextBoxImpl) SetDesiredWidth(desiredWidth int) {
	if t.desiredWidth != desiredWidth {
		t.desiredWidth = desiredWidth
		t.SizeChanged()
	}
}

func (t *TextBoxImpl) Select(list TextSelectionList) {
	t.controller.StoreCaretLocations()
	t.controller.SetSelections(list)
	// Use two scroll tos to try and display all selections (if it fits on screen)
	t.ScrollToRune(t.controller.FirstSelection().First())
	t.ScrollToRune(t.controller.LastSelection().Last())
}

func (t *TextBoxImpl) SelectAll() {
	t.controller.StoreCaretLocations()
	t.controller.SelectAll()
	t.ScrollToRune(t.controller.FirstCaret())
}

func (t *TextBoxImpl) Carets() []int {
	return t.controller.Carets()
}

func (t *TextBoxImpl) RuneIndexAt(point math.Point) (int, bool) {
	for _, child := range ControlsUnder(point, t) {
		line, _ := child.Control.(TextBoxLine)
		if line == nil {
			continue
		}

		point = ParentToChild(point, t.parent, line)
		return line.RuneIndexAt(point), true
	}
	return -1, false
}

func (t *TextBoxImpl) TextAt(start, end int) string {
	return t.controller.TextRange(start, end)
}

func (t *TextBoxImpl) WordAt(runeIndex int) string {
	s, e := t.controller.WordAt(runeIndex)
	return t.controller.TextRange(s, e)
}

func (t *TextBoxImpl) LineIndex(runeIndex int) int {
	return t.controller.LineIndex(runeIndex)
}

func (t *TextBoxImpl) LineStart(line int) int {
	return t.controller.LineStart(line)
}

func (t *TextBoxImpl) LineEnd(line int) int {
	return t.controller.LineEnd(line)
}

func (t *TextBoxImpl) ScrollToLine(index int) {
	t.ListImpl.ScrollTo(index)
}

func (t *TextBoxImpl) ScrollToRune(index int) {
	t.ScrollToLine(t.controller.LineIndex(index))
}

func (t *TextBoxImpl) KeyPress(event KeyboardEvent) bool {
	switch event.Key {
	case KeyLeft:
		switch {
		case event.Modifier.Shift() && event.Modifier.Control():
			t.controller.SelectLeftByWord()
		case event.Modifier.Shift():
			t.controller.SelectLeft()
		case event.Modifier.Alt():
			t.controller.RestorePreviousSelections()
		case !t.controller.Deselect(true):
			if event.Modifier.Control() {
				t.controller.MoveLeftByWord()
			} else {
				t.controller.MoveLeft()
			}
		}
		t.ScrollToRune(t.controller.FirstCaret())
		return true

	case KeyRight:
		switch {
		case event.Modifier.Shift() && event.Modifier.Control():
			t.controller.SelectRightByWord()
		case event.Modifier.Shift():
			t.controller.SelectRight()
		case event.Modifier.Alt():
			t.controller.RestoreNextSelections()
		case !t.controller.Deselect(false):
			if event.Modifier.Control() {
				t.controller.MoveRightByWord()
			} else {
				t.controller.MoveRight()
			}
		}
		t.ScrollToRune(t.controller.LastCaret())
		return true

	case KeyUp:
		switch {
		case event.Modifier.Shift() && event.Modifier.Alt():
			t.controller.AddCaretsUp()
		case event.Modifier.Shift():
			t.controller.SelectUp()
		default:
			t.controller.Deselect(true)
			t.controller.MoveUp()
		}
		t.ScrollToRune(t.controller.FirstCaret())
		return true

	case KeyDown:
		switch {
		case event.Modifier.Shift() && event.Modifier.Alt():
			t.controller.AddCaretsDown()
		case event.Modifier.Shift():
			t.controller.SelectDown()
		default:
			t.controller.Deselect(false)
			t.controller.MoveDown()
		}
		t.ScrollToRune(t.controller.LastCaret())
		return true

	case KeyHome:
		switch {
		case event.Modifier.Shift() && event.Modifier.Control():
			t.controller.SelectFirst()
		case event.Modifier.Control():
			t.controller.MoveFirst()
		case event.Modifier.Shift():
			t.controller.SelectHome()
		default:
			t.controller.Deselect(true)
			t.controller.MoveHome()
		}
		t.ScrollToRune(t.controller.FirstCaret())
		return true

	case KeyEnd:
		switch {
		case event.Modifier.Shift() && event.Modifier.Control():
			t.controller.SelectLast()
		case event.Modifier.Control():
			t.controller.MoveLast()
		case event.Modifier.Shift():
			t.controller.SelectEnd()
		default:
			t.controller.Deselect(false)
			t.controller.MoveEnd()
		}
		t.ScrollToRune(t.controller.LastCaret())
		return true

	case KeyPageUp:
		switch {
		case event.Modifier.Shift():
			for i, c := 0, t.pageLines(); i < c; i++ {
				t.controller.SelectUp()
			}
		default:
			t.controller.Deselect(true)
			for i, c := 0, t.pageLines(); i < c; i++ {
				t.controller.MoveUp()
			}
		}
		t.ScrollToRune(t.controller.FirstCaret())
		return true

	case KeyPageDown:
		switch {
		case event.Modifier.Shift():
			for i, c := 0, t.pageLines(); i < c; i++ {
				t.controller.SelectDown()
			}
		default:
			t.controller.Deselect(false)
			for i, c := 0, t.pageLines(); i < c; i++ {
				t.controller.MoveDown()
			}
		}
		t.ScrollToRune(t.controller.LastCaret())
		return true

	case KeyBackspace:
		t.controller.Backspace()
		return true

	case KeyDelete:
		t.controller.Delete()
		return true

	case KeyEnter:
		if t.multiline {
			t.controller.ReplaceWithNewline()
			return true
		}

	case KeyA:
		if event.Modifier.Control() {
			t.controller.SelectAll()
			return true
		}

	case KeyX:
		fallthrough

	case KeyC:
		if event.Modifier.Control() {
			parts := make([]string, t.controller.SelectionCount())
			for i, _ := range parts {
				parts[i] = t.controller.SelectionText(i)
				if parts[i] == "" {
					// Copy line instead.
					parts[i] = "\n" + t.controller.SelectionLineText(i)
				}
			}
			str := strings.Join(parts, "\n")
			t.driver.SetClipboard(str)

			if event.Key == KeyX {
				t.controller.ReplaceAll("")
			}
			return true
		}

	case KeyV:
		if event.Modifier.Control() {
			str, _ := t.driver.GetClipboard()
			t.controller.ReplaceAll(str)
			t.controller.Deselect(false)
			return true
		}

	case KeyEscape:
		t.controller.ClearSelections()
	}

	return t.ListImpl.KeyPress(event)
}

func (t *TextBoxImpl) KeyStroke(event KeyStrokeEvent) bool {
	if !event.Modifier.Control() && !event.Modifier.Alt() {
		t.controller.ReplaceAllRunes([]rune{event.Character})
		t.controller.Deselect(false)
	}
	t.InputEventHandlerPart.KeyStroke(event)
	return true
}

func (t *TextBoxImpl) Click(event MouseEvent) bool {
	t.InputEventHandlerPart.Click(event)
	return true
}

func (t *TextBoxImpl) DoubleClick(event MouseEvent) bool {
	if p, ok := t.RuneIndexAt(event.Point); ok {
		s, e := t.controller.WordAt(p)
		if event.Modifier&ModControl != 0 {
			t.controller.AddSelection(CreateTextSelection(s, e, false))
		} else {
			t.controller.SetSelection(CreateTextSelection(s, e, false))
		}
	}
	t.InputEventHandlerPart.DoubleClick(event)
	return true
}

func (t *TextBoxImpl) MouseMove(event MouseEvent) {
	t.ListImpl.MouseMove(event)
	if t.selectionDragging {
		if point, ok := t.RuneIndexAt(event.Point); ok {
			t.selectionDrag = CreateTextSelection(t.selectionDrag.From(), point, false)
			t.selectionDragging = true
			t.onRedrawLines.Fire()
		}
	}
}

func (t *TextBoxImpl) CreateLine(driver Driver, styles *StyleDefs, index int) (TextBoxLine, Control) {
	l := &DefaultTextBoxLine{}
	l.Init(l, t, index)
	return l, l
}

// mixins.ListImpl overrides
func (t *TextBoxImpl) PaintSelection(c Canvas, r math.Rect) {}

func (t *TextBoxImpl) PaintMouseOverBackground(c Canvas, r math.Rect) {}

// gxui.AdapterCompliance
type TextBoxAdapter struct {
	DefaultAdapter
	TextBox *TextBoxImpl
}

func (t *TextBoxAdapter) Count() int {
	return math.Max(t.TextBox.controller.LineCount(), 1)
}

func (t *TextBoxAdapter) ItemAt(index int) AdapterItem {
	return index
}

func (t *TextBoxAdapter) ItemIndex(item AdapterItem) int {
	return item.(int)
}

func (t *TextBoxAdapter) Size(styles *StyleDefs) math.Size {
	tb := t.TextBox
	return math.Size{W: tb.desiredWidth, H: tb.font.GlyphMaxSize().H}
}

func (t *TextBoxAdapter) Create(driver Driver, styles *StyleDefs, index int) Control {
	line, container := t.TextBox.parent.CreateLine(driver, styles, index)
	line.OnMouseDown(
		func(ev MouseEvent) {
			t.TextBox.lineMouseDown(line, ev)
		},
	)
	line.OnMouseUp(
		func(ev MouseEvent) {
			t.TextBox.lineMouseUp(line, ev)
		},
	)
	return container
}
