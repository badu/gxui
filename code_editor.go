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

type CodeEditorParent interface {
	TextBoxParent
	CreateSuggestionList() *ListImpl
}

type CodeEditor struct {
	TextBox
	parent             CodeEditorParent
	layers             CodeSyntaxLayers
	suggestionList     *ListImpl
	suggestionProvider CodeSuggestionProvider
	styles             *StyleDefs
	suggestionAdapter  *SuggestionAdapter
	tabWidth           int
}

func (t *CodeEditor) Init(parent CodeEditorParent, driver Driver, styles *StyleDefs) {
	t.parent = parent
	t.driver = driver
	t.styles = styles

	t.tabWidth = int(styles.CodeEditorStyle.Pen.Width)

	t.suggestionAdapter = &SuggestionAdapter{}
	t.suggestionList = t.parent.CreateSuggestionList()
	t.suggestionList.SetAdapter(t.suggestionAdapter)

	t.TextBox.Init(parent, driver, styles, styles.CodeEditorStyle.Font)
	t.controller.OnTextChanged(t.updateSpans)
}

func (t *CodeEditor) ItemSize(styles *StyleDefs) math.Size {
	return math.Size{W: math.MaxSize.W, H: t.font.GlyphMaxSize().H}
}

func (t *CodeEditor) CreateSuggestionList() *ListImpl {
	list := CreateList(t.driver, t.styles)
	list.SetBackgroundBrush(DefaultBrush)
	list.SetBorderPen(DefaultPen)
	return list
}

func (t *CodeEditor) SyntaxLayers() CodeSyntaxLayers {
	return t.layers
}

func (t *CodeEditor) SetSyntaxLayers(layers CodeSyntaxLayers) {
	t.layers = layers
	t.onRedrawLines.Emit()
}

func (t *CodeEditor) TabWidth() int {
	return t.tabWidth
}

func (t *CodeEditor) SetTabWidth(tabWidth int) {
	t.tabWidth = tabWidth
}

func (t *CodeEditor) SuggestionProvider() CodeSuggestionProvider {
	return t.suggestionProvider
}

func (t *CodeEditor) SetSuggestionProvider(provider CodeSuggestionProvider) {
	if t.suggestionProvider != provider {
		t.suggestionProvider = provider
		if t.IsSuggestionListShowing() {
			t.ShowSuggestionList() // Update list
		}
	}
}

func (t *CodeEditor) IsSuggestionListShowing() bool {
	return t.parent.Children().Find(t.suggestionList) != nil
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
	lineOffset := ChildToParent(math.ZeroPoint, line, t.parent)
	target := line.PositionAt(caret).Add(lineOffset)
	childSize := t.suggestionList.DesiredSize(math.ZeroSize, bounds.Size())

	t.suggestionList.Select(t.suggestionList.Adapter().ItemAt(0))
	t.suggestionList.SetSize(childSize)

	child.Layout(childSize.Rect().Offset(target).Intersect(bounds))
}

func (t *CodeEditor) HideSuggestionList() {
	if !t.IsSuggestionListShowing() {
		return
	}

	t.RemoveChild(t.suggestionList)
}

func (t *CodeEditor) Line(idx int) TextBoxLine {
	return FindControl(
		t.ItemControl(idx).(Parent),
		func(c Control) bool {
			_, b := c.(TextBoxLine)
			return b
		},
	).(TextBoxLine)
}

// mixins.ListImpl overrides
func (t *CodeEditor) Click(event MouseEvent) bool {
	t.HideSuggestionList()
	return t.TextBox.Click(event)
}

func (t *CodeEditor) KeyPress(event KeyboardEvent) bool {
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

	return t.TextBox.KeyPress(event)
}

func (t *CodeEditor) KeyStroke(event KeyStrokeEvent) bool {
	consume := t.TextBox.KeyStroke(event)
	if t.IsSuggestionListShowing() {
		t.SortSuggestionList()
	}

	return consume
}

// mixins.TextBox overrides
func (t *CodeEditor) CreateLine(driver Driver, styles *StyleDefs, index int) (TextBoxLine, Control) {
	lineNumber := CreateLabel(driver, styles)

	lineNumber.SetText(fmt.Sprintf("%d", index+1)) // Displayed lines start at 1

	line := &CodeEditorLine{}
	line.Init(line, t, index)

	layout := CreateLinearLayout(driver, styles)
	layout.SetDirection(LeftToRight)
	layout.AddChild(lineNumber)
	layout.AddChild(line)

	return line, layout
}

func (t *CodeEditor) updateSpans(edits []TextBoxEdit) {
	runeCount := len(t.controller.TextRunes())
	for _, layer := range t.layers {
		layer.UpdateSpans(runeCount, edits)
	}
}
