// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"fmt"
	"github.com/badu/gxui"
	"github.com/badu/gxui/math"
	"strings"
)

type CodeEditor struct {
	TextBox
	outer              gxui.CodeEditorOuter
	layers             gxui.CodeSyntaxLayers
	suggestionAdapter  *SuggestionAdapter
	suggestionList     gxui.List
	suggestionProvider gxui.CodeSuggestionProvider
	tabWidth           int
	theme              gxui.Theme
}

func (t *CodeEditor) Init(outer gxui.CodeEditorOuter, driver gxui.Driver, theme gxui.Theme, font gxui.Font) {
	t.outer = outer
	t.tabWidth = 2
	t.theme = theme

	t.suggestionAdapter = &SuggestionAdapter{}
	t.suggestionList = t.outer.CreateSuggestionList()
	t.suggestionList.SetAdapter(t.suggestionAdapter)

	t.TextBox.Init(outer, driver, theme, font)
	t.controller.OnTextChanged(t.updateSpans)
}

func (t *CodeEditor) ItemSize(theme gxui.Theme) math.Size {
	return math.Size{W: math.MaxSize.W, H: t.font.GlyphMaxSize().H}
}

func (t *CodeEditor) CreateSuggestionList() gxui.List {
	list := t.theme.CreateList()
	list.SetBackgroundBrush(gxui.DefaultBrush)
	list.SetBorderPen(gxui.DefaultPen)
	return list
}

func (t *CodeEditor) SyntaxLayers() gxui.CodeSyntaxLayers {
	return t.layers
}

func (t *CodeEditor) SetSyntaxLayers(layers gxui.CodeSyntaxLayers) {
	t.layers = layers
	t.onRedrawLines.Fire()
}

func (t *CodeEditor) TabWidth() int {
	return t.tabWidth
}

func (t *CodeEditor) SetTabWidth(tabWidth int) {
	t.tabWidth = tabWidth
}

func (t *CodeEditor) SuggestionProvider() gxui.CodeSuggestionProvider {
	return t.suggestionProvider
}

func (t *CodeEditor) SetSuggestionProvider(provider gxui.CodeSuggestionProvider) {
	if t.suggestionProvider != provider {
		t.suggestionProvider = provider
		if t.IsSuggestionListShowing() {
			t.ShowSuggestionList() // Update list
		}
	}
}

func (t *CodeEditor) IsSuggestionListShowing() bool {
	return t.outer.Children().Find(t.suggestionList) != nil
}

func (t *CodeEditor) SortSuggestionList() {
	caret := t.controller.LastCaret()
	partial := t.controller.TextRange(t.controller.WordAt(caret))
	t.suggestionAdapter.Sort(partial)
}

func (t *CodeEditor) ShowSuggestionList() {
	if t.suggestionProvider == nil || t.IsSuggestionListShowing() {
		return
	}

	caret := t.controller.LastCaret()
	word, _ := t.controller.WordAt(caret)

	suggestions := t.suggestionProvider.SuggestionsAt(word)
	if len(suggestions) == 0 {
		t.HideSuggestionList()
		return
	}

	t.suggestionAdapter.SetSuggestions(suggestions)
	t.SortSuggestionList()
	child := t.AddChild(t.suggestionList)

	// Position the suggestion list below the last caret
	lineIdx := t.controller.LineIndex(caret)
	// TODO: What if the last caret is not visible?
	bounds := t.Size().Rect().Contract(t.Padding())
	line := t.Line(lineIdx)
	lineOffset := gxui.ChildToParent(math.ZeroPoint, line, t.outer)
	target := line.PositionAt(caret).Add(lineOffset)
	cs := t.suggestionList.DesiredSize(math.ZeroSize, bounds.Size())
	t.suggestionList.Select(t.suggestionList.Adapter().ItemAt(0))
	t.suggestionList.SetSize(cs)
	child.Layout(cs.Rect().Offset(target).Intersect(bounds))
}

func (t *CodeEditor) HideSuggestionList() {
	if t.IsSuggestionListShowing() {
		t.RemoveChild(t.suggestionList)
	}
}

func (t *CodeEditor) Line(idx int) gxui.TextBoxLine {
	return gxui.FindControl(
		t.ItemControl(idx).(gxui.Parent),
		func(c gxui.Control) bool {
			_, b := c.(gxui.TextBoxLine)
			return b
		},
	).(gxui.TextBoxLine)
}

// mixins.List overrides
func (t *CodeEditor) Click(event gxui.MouseEvent) bool {
	t.HideSuggestionList()
	return t.TextBox.Click(event)
}

func (t *CodeEditor) KeyPress(event gxui.KeyboardEvent) bool {
	switch event.Key {
	case gxui.KeyTab:
		replace := true
		for _, selection := range t.controller.Selections() {
			start, end := selection.Range()
			if t.controller.LineIndex(start) != t.controller.LineIndex(end) {
				replace = false
				break
			}
		}
		switch {
		case replace:
			t.controller.ReplaceAll(strings.Repeat(" ", t.tabWidth))
			t.controller.Deselect(false)
		case event.Modifier.Shift():
			t.controller.UnindentSelection(t.tabWidth)
		default:
			t.controller.IndentSelection(t.tabWidth)
		}
		return true

	case gxui.KeySpace:
		if event.Modifier.Control() {
			t.ShowSuggestionList()
			return false
		}

	case gxui.KeyUp:
		fallthrough

	case gxui.KeyDown:
		if t.IsSuggestionListShowing() {
			return t.suggestionList.KeyPress(event)
		}

	case gxui.KeyLeft:
		t.HideSuggestionList()

	case gxui.KeyRight:
		t.HideSuggestionList()

	case gxui.KeyEnter:
		controller := t.controller
		if t.IsSuggestionListShowing() {
			text := t.suggestionAdapter.Suggestion(t.suggestionList.Selected()).Code()
			start, end := controller.WordAt(t.controller.LastCaret())
			controller.SetSelection(gxui.CreateTextSelection(start, end, false))
			controller.ReplaceAll(text)
			controller.Deselect(false)
			t.HideSuggestionList()
		} else {
			t.controller.ReplaceWithNewlineKeepIndent()
		}
		return true

	case gxui.KeyEscape:
		if t.IsSuggestionListShowing() {
			t.HideSuggestionList()
			return true
		}
	}

	return t.TextBox.KeyPress(event)
}

func (t *CodeEditor) KeyStroke(event gxui.KeyStrokeEvent) bool {
	consume := t.TextBox.KeyStroke(event)
	if t.IsSuggestionListShowing() {
		t.SortSuggestionList()
	}

	return consume
}

// mixins.TextBox overrides
func (t *CodeEditor) CreateLine(theme gxui.Theme, index int) (gxui.TextBoxLine, gxui.Control) {
	lineNumber := theme.CreateLabel()
	lineNumber.SetText(fmt.Sprintf("%d", index+1)) // Displayed lines start at 1

	line := &CodeEditorLine{}
	line.Init(line, theme, t, index)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(gxui.LeftToRight)
	layout.AddChild(lineNumber)
	layout.AddChild(line)

	return line, layout
}

func (t *CodeEditor) updateSpans(edits []gxui.TextBoxEdit) {
	runeCount := len(t.controller.TextRunes())
	for _, layer := range t.layers {
		layer.UpdateSpans(runeCount, edits)
	}
}
