// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"strings"

	"github.com/badu/gxui/math"
)

type LineWidth struct {
	Tabs  int
	Chars int
}

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
	suggestionProvider CodeSuggestionProvider
	suggestionList     *ListImpl
	styles             *StyleDefs
	suggestionAdapter  *SuggestionAdapter
	hiddenLines        map[int]struct{}
	layers             CodeSyntaxLayers
	lineWidths         []LineWidth
	tabWidth           int
}

func (e *CodeEditor) Init(parent CodeEditorParent, driver Driver, styles *StyleDefs) {
	e.parent = parent
	e.driver = driver
	e.styles = styles

	e.tabWidth = int(styles.CodeEditorStyle.Pen.Width)
	e.hiddenLines = map[int]struct{}{}

	e.suggestionAdapter = &SuggestionAdapter{}
	e.suggestionList = e.parent.CreateSuggestionList()
	e.suggestionList.SetAdapter(e.suggestionAdapter)

	e.TextBox.Init(parent, driver, styles, styles.CodeEditorStyle.Font)
	e.TextBox.horizontalScroll.Forget()
	e.TextBox.horizontalScroll = e.TextBox.horizontalScrollbar.OnScroll(
		func(from, to int) {
			e.SetHorizontalOffset(from)
		},
	)

	e.controller.OnTextChanged(e.updateSpans)
}

func (e *CodeEditor) ItemSize(styles *StyleDefs) math.Size {
	return math.Size{W: math.MaxSize.W, H: e.font.GlyphMaxSize().H}
}

func (e *CodeEditor) CreateSuggestionList() *ListImpl {
	list := CreateList(e.driver, e.styles)
	list.SetBackgroundBrush(DefaultBrush)
	list.SetBorderPen(DefaultPen)
	return list
}

func (e *CodeEditor) SyntaxLayers() CodeSyntaxLayers {
	if len(e.layers) == 0 {
		e.layers = append(e.layers, CreateCodeSyntaxLayer())
	}
	return e.layers
}

func (e *CodeEditor) SetSyntaxLayers(layers CodeSyntaxLayers) {
	e.layers = layers
	e.onRedrawLines.Emit()
}

func (e *CodeEditor) TabWidth() int {
	return e.tabWidth
}

func (e *CodeEditor) SetTabWidth(tabWidth int) {
	e.tabWidth = tabWidth
}

func (e *CodeEditor) SuggestionProvider() CodeSuggestionProvider {
	return e.suggestionProvider
}

func (e *CodeEditor) SetSuggestionProvider(provider CodeSuggestionProvider) {
	if e.suggestionProvider == provider {
		return
	}

	e.suggestionProvider = provider
	if e.IsSuggestionListShowing() {
		e.ShowSuggestionList() // Update list
	}
}

func (e *CodeEditor) IsSuggestionListShowing() bool {
	return e.parent.Children().Find(e.suggestionList) != nil
}

func (e *CodeEditor) SortSuggestionList() {
	caret := e.controller.LastCaret()
	partial := e.controller.TextRange(e.controller.WordAt(caret))
	e.suggestionAdapter.Sort(partial)
}

func (e *CodeEditor) ShowSuggestionList() {
	if e.suggestionProvider == nil || e.IsSuggestionListShowing() {
		return
	}

	caret := e.controller.LastCaret()
	word, _ := e.controller.WordAt(caret)

	suggestions := e.suggestionProvider.SuggestionsAt(word)
	if len(suggestions) == 0 {
		e.HideSuggestionList()
		return
	}

	e.suggestionAdapter.SetSuggestions(suggestions)
	e.SortSuggestionList()
	child := e.AddChild(e.suggestionList)

	// Position the suggestion list below the last caret
	lineIdx := e.controller.LineIndex(caret)
	// TODO: What if the last caret is not visible?
	bounds := e.Size().Rect().Contract(e.Padding())
	line := e.Line(lineIdx)
	lineOffset := ChildToParent(math.ZeroPoint, line, e.parent)
	target := line.PositionAt(caret).Add(lineOffset)
	childSize := e.suggestionList.DesiredSize(math.ZeroSize, bounds.Size())

	e.suggestionList.Select(e.suggestionList.Adapter().ItemAt(0))
	e.suggestionList.SetSize(childSize)

	child.Layout(childSize.Rect().Offset(target).Intersect(bounds))
}

func (e *CodeEditor) HideSuggestionList() {
	if !e.IsSuggestionListShowing() {
		return
	}

	e.RemoveChild(e.suggestionList)
}

func (e *CodeEditor) Line(idx int) TextBoxLine {
	return FindControl(
		e.ItemControl(idx).(Parent),
		func(c Control) bool {
			_, b := c.(TextBoxLine)
			return b
		},
	).(TextBoxLine)
}

// mixins.ListImpl overrides
func (e *CodeEditor) Click(event MouseEvent) bool {
	e.HideSuggestionList()
	return e.TextBox.Click(event)
}

func (e *CodeEditor) KeyPress(event KeyboardEvent) bool {
	switch event.Key {
	case KeyTab:
		replace := true
		for _, selection := range e.controller.Selections() {
			start, end := selection.Range()
			if e.controller.LineIndex(start) != e.controller.LineIndex(end) {
				replace = false
				break
			}
		}

		switch {
		case replace:
			e.controller.ReplaceAll(strings.Repeat(" ", e.tabWidth))
			e.controller.Deselect(false)
		case event.Modifier.Shift():
			e.controller.UnindentSelection(e.tabWidth)
		default:
			e.controller.IndentSelection(e.tabWidth)
		}

		return true

	case KeySpace:
		if event.Modifier.Control() {
			e.ShowSuggestionList()
			return false
		}

	case KeyUp:
		fallthrough

	case KeyDown:
		if e.IsSuggestionListShowing() {
			return e.suggestionList.KeyPress(event)
		}

	case KeyLeft:
		e.HideSuggestionList()

	case KeyRight:
		e.HideSuggestionList()

	case KeyEnter:
		controller := e.controller
		if e.IsSuggestionListShowing() {
			text := e.suggestionAdapter.Suggestion(e.suggestionList.Selected()).Code()
			start, end := controller.WordAt(e.controller.LastCaret())
			controller.SetSelection(CreateTextSelection(start, end, false))
			controller.ReplaceAll(text)
			controller.Deselect(false)
			e.HideSuggestionList()
		} else {
			e.controller.ReplaceWithNewlineKeepIndent()
		}

		return true

	case KeyEscape:
		if e.IsSuggestionListShowing() {
			e.HideSuggestionList()
			return true
		}
	}

	return e.TextBox.KeyPress(event)
}

func (e *CodeEditor) KeyStroke(event KeyStrokeEvent) bool {
	consume := e.TextBox.KeyStroke(event)
	if e.IsSuggestionListShowing() {
		e.SortSuggestionList()
	}

	return consume
}

// mixins.TextBox overrides
func (e *CodeEditor) CreateLine(driver Driver, styles *StyleDefs, index int) (TextBoxLine, Control) {
	lineNumber := CreateLabel(driver, styles)

	lineNumber.SetText(fmt.Sprintf("%d", index+1)) // Displayed lines start at 1

	line := &CodeEditorLine{}
	line.Init(line, e, index)

	foldButton := CreateButton(driver, styles)
	foldButton.SetMargin(math.Spacing{L: 0, T: 0, R: 0, B: 0})
	foldButton.SetVisible(false)

	layout := CreateLinearLayout(driver, styles)
	layout.SetDirection(LeftToRight)
	layout.AddChild(lineNumber)
	layout.AddChild(foldButton)
	layout.AddChild(line)

	if _, ok := e.hiddenLines[index]; ok {
		layout.SetVisible(false)
	}

	return line, layout
}

func (e *CodeEditor) updateSpans(edits []TextBoxEdit) {
	runeCount := len(e.controller.TextRunes())
	for _, layer := range e.layers {
		layer.UpdateSpans(runeCount, edits)
	}
}

func (e *CodeEditor) RevealLines(from int, to int) {
	for i := from; i <= to; i++ {
		delete(e.hiddenLines, i)
		e.ChangeHiddenCount(-1)
	}
	e.LayoutChildren()
}

func (e *CodeEditor) HideLines(from int, to int) {
	for i := from; i <= to; i++ {
		e.hiddenLines[i] = struct{}{}
		e.ChangeHiddenCount(1)
		ctrl := e.ItemControl(i)
		if ctrl != nil {
			ctrl.SetVisible(false)
		}
	}
	e.LayoutChildren()
}

func (e *CodeEditor) SetHorizontalOffset(offset int) {
	if e.horizontalOffset == offset {
		return
	}
	e.updateHorizScrollLimit()
	e.updateChildOffsets(e, offset)
	e.horizontalScrollbar.SetScrollPosition(offset, offset+e.Size().W)
	e.horizontalOffset = offset
	e.LayoutChildren()
}

func (e *CodeEditor) updateHorizScrollLimit() {
	maxWidth := e.MaxLineWidth()
	size := e.Size().Contract(e.parent.Padding())
	maxScroll := math.Max(maxWidth-size.W, 0)
	math.Clamp(e.horizontalOffset, 0, maxScroll)
	e.horizontalScrollbar.SetScrollLimit(maxWidth)
}

func (e *CodeEditor) SetLineWidths(widths []LineWidth) {
	e.lineWidths = widths
}

func (e *CodeEditor) MaxLineWidth() int {
	if len(e.lineWidths) == 0 {
		return 0
	}

	maxWidth := LineWidth{}
	for index, lineWidth := range e.lineWidths {

		width := lineWidth.Tabs*e.TabWidth() + lineWidth.Chars

		if width <= maxWidth.Chars {
			continue
		}

		maxWidth.Tabs = index
		maxWidth.Chars = width
	}

	line, _ := e.CreateLine(e.driver, e.styles, maxWidth.Tabs)
	lineEnd := e.controller.LineEnd(maxWidth.Tabs)
	lastPos := line.PositionAt(lineEnd)

	return e.lineWidthOffset() + lastPos.X
}
