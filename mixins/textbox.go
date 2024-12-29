// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui/mixins/base"
	"strings"

	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
)

type TextBoxLine interface {
	gxui.Control
	RuneIndexAt(math.Point) int
	PositionAt(int) math.Point
}

type TextBoxOuter interface {
	ListOuter
	CreateLine(theme gxui.Theme, index int) (line TextBoxLine, container gxui.Control)
}

type TextBox struct {
	List
	gxui.AdapterBase
	base.FocusablePart
	outer             TextBoxOuter
	driver            gxui.Driver
	font              gxui.Font
	textColor         gxui.Color
	onRedrawLines     gxui.Event
	multiline         bool
	controller        *gxui.TextBoxController
	adapter           *TextBoxAdapter
	selectionDragging bool
	selectionDrag     gxui.TextSelection
	desiredWidth      int
}

func (t *TextBox) lineMouseDown(line TextBoxLine, event gxui.MouseEvent) {
	if event.Button == gxui.MouseButtonLeft {
		p := line.RuneIndexAt(event.Point)
		t.selectionDragging = true
		t.selectionDrag = gxui.CreateTextSelection(p, p, false)
		if !event.Modifier.Control() {
			t.controller.SetCaret(p)
		}
	}
}

func (t *TextBox) lineMouseUp(line TextBoxLine, event gxui.MouseEvent) {
	if event.Button == gxui.MouseButtonLeft {
		t.selectionDragging = false
		if !event.Modifier.Control() {
			t.controller.SetSelection(t.selectionDrag)
		} else {
			t.controller.AddSelection(t.selectionDrag)
		}
	}
}

func (t *TextBox) Init(outer TextBoxOuter, driver gxui.Driver, theme gxui.Theme, font gxui.Font) {
	t.List.Init(outer, theme)
	t.FocusablePart.Init(outer)
	t.outer = outer
	t.driver = driver
	t.font = font
	t.onRedrawLines = gxui.CreateEvent(func() {})
	t.controller = gxui.CreateTextBoxController()
	t.adapter = &TextBoxAdapter{TextBox: t}
	t.desiredWidth = 100
	t.SetScrollBarEnabled(false) // Defaults to single line
	t.OnGainedFocus(func() { t.onRedrawLines.Fire() })
	t.OnLostFocus(func() { t.onRedrawLines.Fire() })
	t.controller.OnTextChanged(func([]gxui.TextBoxEdit) {
		t.onRedrawLines.Fire()
		t.List.DataChanged(false)
	})
	t.controller.OnSelectionChanged(func() {
		t.onRedrawLines.Fire()
	})

	t.List.SetAdapter(t.adapter)
}

func (t *TextBox) textRect() math.Rect {
	return t.outer.Size().Rect().Contract(t.Padding())
}

func (t *TextBox) pageLines() int {
	return (t.outer.Size().H - t.outer.Padding().H()) / t.MajorAxisItemSize()
}

func (t *TextBox) OnRedrawLines(callback func()) gxui.EventSubscription {
	return t.onRedrawLines.Listen(callback)
}

func (t *TextBox) OnSelectionChanged(callback func()) gxui.EventSubscription {
	return t.controller.OnSelectionChanged(callback)
}

func (t *TextBox) OnTextChanged(callback func(lines []gxui.TextBoxEdit)) gxui.EventSubscription {
	return t.controller.OnTextChanged(callback)
}

func (t *TextBox) Runes() []rune {
	return t.controller.TextRunes()
}

func (t *TextBox) Text() string {
	return t.controller.Text()
}

func (t *TextBox) SetText(text string) {
	t.controller.SetText(text)
	t.outer.Relayout()
}

func (t *TextBox) TextColor() gxui.Color {
	return t.textColor
}

func (t *TextBox) SetTextColor(color gxui.Color) {
	t.textColor = color
	t.Relayout()
}

func (t *TextBox) Font() gxui.Font {
	return t.font
}

func (t *TextBox) SetFont(font gxui.Font) {
	if t.font != font {
		t.font = font
		t.Relayout()
	}
}

func (t *TextBox) Multiline() bool {
	return t.multiline
}

func (t *TextBox) SetMultiline(multiline bool) {
	if t.multiline != multiline {
		t.multiline = multiline
		t.SetScrollBarEnabled(multiline)
		t.outer.Relayout()
	}
}

func (t *TextBox) DesiredWidth() int {
	return t.desiredWidth
}

func (t *TextBox) SetDesiredWidth(desiredWidth int) {
	if t.desiredWidth != desiredWidth {
		t.desiredWidth = desiredWidth
		t.SizeChanged()
	}
}

func (t *TextBox) Select(list gxui.TextSelectionList) {
	t.controller.StoreCaretLocations()
	t.controller.SetSelections(list)
	// Use two scroll tos to try and display all selections (if it fits on screen)
	t.ScrollToRune(t.controller.FirstSelection().First())
	t.ScrollToRune(t.controller.LastSelection().Last())
}

func (t *TextBox) SelectAll() {
	t.controller.StoreCaretLocations()
	t.controller.SelectAll()
	t.ScrollToRune(t.controller.FirstCaret())
}

func (t *TextBox) Carets() []int {
	return t.controller.Carets()
}

func (t *TextBox) RuneIndexAt(point math.Point) (int, bool) {
	for _, child := range gxui.ControlsUnder(point, t) {
		line, _ := child.Control.(TextBoxLine)
		if line == nil {
			continue
		}

		point = gxui.ParentToChild(point, t.outer, line)
		return line.RuneIndexAt(point), true
	}
	return -1, false
}

func (t *TextBox) TextAt(start, end int) string {
	return t.controller.TextRange(start, end)
}

func (t *TextBox) WordAt(runeIndex int) string {
	s, e := t.controller.WordAt(runeIndex)
	return t.controller.TextRange(s, e)
}

func (t *TextBox) LineIndex(runeIndex int) int {
	return t.controller.LineIndex(runeIndex)
}

func (t *TextBox) LineStart(line int) int {
	return t.controller.LineStart(line)
}

func (t *TextBox) LineEnd(line int) int {
	return t.controller.LineEnd(line)
}

func (t *TextBox) ScrollToLine(index int) {
	t.List.ScrollTo(index)
}

func (t *TextBox) ScrollToRune(index int) {
	t.ScrollToLine(t.controller.LineIndex(index))
}

func (t *TextBox) KeyPress(event gxui.KeyboardEvent) bool {
	switch event.Key {
	case gxui.KeyLeft:
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

	case gxui.KeyRight:
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

	case gxui.KeyUp:
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

	case gxui.KeyDown:
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

	case gxui.KeyHome:
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

	case gxui.KeyEnd:
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

	case gxui.KeyPageUp:
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

	case gxui.KeyPageDown:
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

	case gxui.KeyBackspace:
		t.controller.Backspace()
		return true

	case gxui.KeyDelete:
		t.controller.Delete()
		return true

	case gxui.KeyEnter:
		if t.multiline {
			t.controller.ReplaceWithNewline()
			return true
		}

	case gxui.KeyA:
		if event.Modifier.Control() {
			t.controller.SelectAll()
			return true
		}

	case gxui.KeyX:
		fallthrough

	case gxui.KeyC:
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

			if event.Key == gxui.KeyX {
				t.controller.ReplaceAll("")
			}
			return true
		}

	case gxui.KeyV:
		if event.Modifier.Control() {
			str, _ := t.driver.GetClipboard()
			t.controller.ReplaceAll(str)
			t.controller.Deselect(false)
			return true
		}

	case gxui.KeyEscape:
		t.controller.ClearSelections()
	}

	return t.List.KeyPress(event)
}

func (t *TextBox) KeyStroke(event gxui.KeyStrokeEvent) bool {
	if !event.Modifier.Control() && !event.Modifier.Alt() {
		t.controller.ReplaceAllRunes([]rune{event.Character})
		t.controller.Deselect(false)
	}
	t.InputEventHandlerPart.KeyStroke(event)
	return true
}

func (t *TextBox) Click(event gxui.MouseEvent) bool {
	t.InputEventHandlerPart.Click(event)
	return true
}

func (t *TextBox) DoubleClick(event gxui.MouseEvent) bool {
	if p, ok := t.RuneIndexAt(event.Point); ok {
		s, e := t.controller.WordAt(p)
		if event.Modifier&gxui.ModControl != 0 {
			t.controller.AddSelection(gxui.CreateTextSelection(s, e, false))
		} else {
			t.controller.SetSelection(gxui.CreateTextSelection(s, e, false))
		}
	}
	t.InputEventHandlerPart.DoubleClick(event)
	return true
}

func (t *TextBox) MouseMove(event gxui.MouseEvent) {
	t.List.MouseMove(event)
	if t.selectionDragging {
		if point, ok := t.RuneIndexAt(event.Point); ok {
			t.selectionDrag = gxui.CreateTextSelection(t.selectionDrag.From(), point, false)
			t.selectionDragging = true
			t.onRedrawLines.Fire()
		}
	}
}

func (t *TextBox) CreateLine(theme gxui.Theme, index int) (TextBoxLine, gxui.Control) {
	l := &DefaultTextBoxLine{}
	l.Init(l, theme, t, index)
	return l, l
}

// mixins.List overrides
func (t *TextBox) PaintSelection(c gxui.Canvas, r math.Rect) {}

func (t *TextBox) PaintMouseOverBackground(c gxui.Canvas, r math.Rect) {}

// gxui.AdapterCompliance
type TextBoxAdapter struct {
	gxui.DefaultAdapter
	TextBox *TextBox
}

func (t *TextBoxAdapter) Count() int {
	return math.Max(t.TextBox.controller.LineCount(), 1)
}

func (t *TextBoxAdapter) ItemAt(index int) gxui.AdapterItem {
	return index
}

func (t *TextBoxAdapter) ItemIndex(item gxui.AdapterItem) int {
	return item.(int)
}

func (t *TextBoxAdapter) Size(theme gxui.Theme) math.Size {
	tb := t.TextBox
	return math.Size{W: tb.desiredWidth, H: tb.font.GlyphMaxSize().H}
}

func (t *TextBoxAdapter) Create(theme gxui.Theme, index int) gxui.Control {
	line, container := t.TextBox.outer.CreateLine(theme, index)
	line.OnMouseDown(func(ev gxui.MouseEvent) {
		t.TextBox.lineMouseDown(line, ev)
	})
	line.OnMouseUp(func(ev gxui.MouseEvent) {
		t.TextBox.lineMouseUp(line, ev)
	})
	return container
}
