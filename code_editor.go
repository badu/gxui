// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"github.com/badu/gxui/math"
	"strings"
)

type CodeSuggestion interface {
	Name() string
	Code() string
}

type CodeSuggestionProvider interface {
	SuggestionsAt(runeIndex int) []CodeSuggestion
}

type CodeEditor interface {
	TextBox
	SyntaxLayers() CodeSyntaxLayers
	SetSyntaxLayers(CodeSyntaxLayers)
	TabWidth() int
	SetTabWidth(int)
	SuggestionProvider() CodeSuggestionProvider
	SetSuggestionProvider(CodeSuggestionProvider)
	ShowSuggestionList()
	HideSuggestionList()
}

type CodeEditorOuter interface {
	ParentTextBox
	CreateSuggestionList() List
}

type CodeEditorImpl struct {
	TextBoxImpl
	outer              CodeEditorOuter
	layers             CodeSyntaxLayers
	suggestionAdapter  *SuggestionAdapter
	suggestionList     List
	suggestionProvider CodeSuggestionProvider
	tabWidth           int
	theme              App
}

func (t *CodeEditorImpl) Init(outer CodeEditorOuter, driver Driver, theme App, font Font) {
	t.outer = outer
	t.tabWidth = 2
	t.theme = theme

	t.suggestionAdapter = &SuggestionAdapter{}
	t.suggestionList = t.outer.CreateSuggestionList()
	t.suggestionList.SetAdapter(t.suggestionAdapter)

	t.TextBoxImpl.Init(outer, driver, theme, font)
	t.controller.OnTextChanged(t.updateSpans)
}

func (t *CodeEditorImpl) ItemSize(theme App) math.Size {
	return math.Size{W: math.MaxSize.W, H: t.font.GlyphMaxSize().H}
}

func (t *CodeEditorImpl) CreateSuggestionList() List {
	list := t.theme.CreateList()
	list.SetBackgroundBrush(DefaultBrush)
	list.SetBorderPen(DefaultPen)
	return list
}

func (t *CodeEditorImpl) SyntaxLayers() CodeSyntaxLayers {
	return t.layers
}

func (t *CodeEditorImpl) SetSyntaxLayers(layers CodeSyntaxLayers) {
	t.layers = layers
	t.onRedrawLines.Fire()
}

func (t *CodeEditorImpl) TabWidth() int {
	return t.tabWidth
}

func (t *CodeEditorImpl) SetTabWidth(tabWidth int) {
	t.tabWidth = tabWidth
}

func (t *CodeEditorImpl) SuggestionProvider() CodeSuggestionProvider {
	return t.suggestionProvider
}

func (t *CodeEditorImpl) SetSuggestionProvider(provider CodeSuggestionProvider) {
	if t.suggestionProvider != provider {
		t.suggestionProvider = provider
		if t.IsSuggestionListShowing() {
			t.ShowSuggestionList() // Update list
		}
	}
}

func (t *CodeEditorImpl) IsSuggestionListShowing() bool {
	return t.outer.Children().Find(t.suggestionList) != nil
}

func (t *CodeEditorImpl) SortSuggestionList() {
	caret := t.controller.LastCaret()
	partial := t.controller.TextRange(t.controller.WordAt(caret))
	t.suggestionAdapter.Sort(partial)
}

func (t *CodeEditorImpl) ShowSuggestionList() {
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
	lineOffset := ChildToParent(math.ZeroPoint, line, t.outer)
	target := line.PositionAt(caret).Add(lineOffset)
	cs := t.suggestionList.DesiredSize(math.ZeroSize, bounds.Size())
	t.suggestionList.Select(t.suggestionList.Adapter().ItemAt(0))
	t.suggestionList.SetSize(cs)
	child.Layout(cs.Rect().Offset(target).Intersect(bounds))
}

func (t *CodeEditorImpl) HideSuggestionList() {
	if t.IsSuggestionListShowing() {
		t.RemoveChild(t.suggestionList)
	}
}

func (t *CodeEditorImpl) Line(idx int) TextBoxLine {
	return FindControl(
		t.ItemControl(idx).(Parent),
		func(c Control) bool {
			_, b := c.(TextBoxLine)
			return b
		},
	).(TextBoxLine)
}

// mixins.ListImpl overrides
func (t *CodeEditorImpl) Click(event MouseEvent) bool {
	t.HideSuggestionList()
	return t.TextBoxImpl.Click(event)
}

func (t *CodeEditorImpl) KeyPress(event KeyboardEvent) bool {
	switch event.Key {
	case KeyTab:
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

	case KeySpace:
		if event.Modifier.Control() {
			t.ShowSuggestionList()
			return false
		}

	case KeyUp:
		fallthrough

	case KeyDown:
		if t.IsSuggestionListShowing() {
			return t.suggestionList.KeyPress(event)
		}

	case KeyLeft:
		t.HideSuggestionList()

	case KeyRight:
		t.HideSuggestionList()

	case KeyEnter:
		controller := t.controller
		if t.IsSuggestionListShowing() {
			text := t.suggestionAdapter.Suggestion(t.suggestionList.Selected()).Code()
			start, end := controller.WordAt(t.controller.LastCaret())
			controller.SetSelection(CreateTextSelection(start, end, false))
			controller.ReplaceAll(text)
			controller.Deselect(false)
			t.HideSuggestionList()
		} else {
			t.controller.ReplaceWithNewlineKeepIndent()
		}
		return true

	case KeyEscape:
		if t.IsSuggestionListShowing() {
			t.HideSuggestionList()
			return true
		}
	}

	return t.TextBoxImpl.KeyPress(event)
}

func (t *CodeEditorImpl) KeyStroke(event KeyStrokeEvent) bool {
	consume := t.TextBoxImpl.KeyStroke(event)
	if t.IsSuggestionListShowing() {
		t.SortSuggestionList()
	}

	return consume
}

// mixins.TextBoxImpl overrides
func (t *CodeEditorImpl) CreateLine(theme App, index int) (TextBoxLine, Control) {
	lineNumber := theme.CreateLabel()
	lineNumber.SetText(fmt.Sprintf("%d", index+1)) // Displayed lines start at 1

	line := &CodeEditorLine{}
	line.Init(line, theme, t, index)

	layout := theme.CreateLinearLayout()
	layout.SetDirection(LeftToRight)
	layout.AddChild(lineNumber)
	layout.AddChild(line)

	return line, layout
}

func (t *CodeEditorImpl) updateSpans(edits []TextBoxEdit) {
	runeCount := len(t.controller.TextRunes())
	for _, layer := range t.layers {
		layer.UpdateSpans(runeCount, edits)
	}
}
