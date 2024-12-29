// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/interval"
	"github.com/badu/gxui/math"
	"sort"
	"strings"
	"unicode"
)

type TextBoxEdit struct {
	At    int
	Delta int
}

type TextBoxController struct {
	onSelectionChanged          Event
	onTextChanged               Event
	text                        []rune
	lineStarts                  []int
	lineEnds                    []int
	selections                  TextSelectionList
	locationHistory             [][]int
	locationHistoryIndex        int
	storeCaretLocationsNextEdit bool
}

func CreateTextBoxController() *TextBoxController {
	result := &TextBoxController{
		onSelectionChanged: CreateEvent(func() {}),
		onTextChanged:      CreateEvent(func([]TextBoxEdit) {}),
	}
	result.selections = TextSelectionList{TextSelection{}}
	return result
}

func (t *TextBoxController) textEdited(edits []TextBoxEdit) {
	t.updateSelectionsForEdits(edits)
	t.onTextChanged.Fire(edits)
}

func (t *TextBoxController) updateSelectionsForEdits(edits []TextBoxEdit) {
	minLen := 0
	maxLen := len(t.text)
	selections := TextSelectionList{}
	for _, selection := range t.selections {
		for _, edit := range edits {
			at := edit.At
			delta := edit.Delta
			if selection.start > at {
				selection.start += delta
			}
			if selection.end >= at {
				selection.end += delta
			}
		}
		if selection.end < selection.start {
			selection.end = selection.start
		}
		selection.start = math.Clamp(selection.start, minLen, maxLen)
		selection.end = math.Clamp(selection.end, minLen, maxLen)
		interval.Merge(&selections, selection)
	}
	t.selections = selections
}

func (t *TextBoxController) setTextRunesNoEvent(text []rune) {
	t.text = text
	t.lineStarts = t.lineStarts[:0]
	t.lineEnds = t.lineEnds[:0]

	t.lineStarts = append(t.lineStarts, 0)
	for index, curRune := range text {
		if curRune == '\n' {
			t.lineEnds = append(t.lineEnds, index)
			t.lineStarts = append(t.lineStarts, index+1)
		}
	}
	t.lineEnds = append(t.lineEnds, len(text))
}

func (t *TextBoxController) maybeStoreCaretLocations() {
	if t.storeCaretLocationsNextEdit {
		t.StoreCaretLocations()
		t.storeCaretLocationsNextEdit = false
	}
}

func (t *TextBoxController) StoreCaretLocations() {
	if t.locationHistoryIndex < len(t.locationHistory) {
		t.locationHistory = t.locationHistory[:t.locationHistoryIndex]
	}
	t.locationHistory = append(t.locationHistory, t.Carets())
	t.locationHistoryIndex = len(t.locationHistory)
}

func (t *TextBoxController) OnSelectionChanged(callback func()) EventSubscription {
	return t.onSelectionChanged.Listen(callback)
}

func (t *TextBoxController) OnTextChanged(callback func([]TextBoxEdit)) EventSubscription {
	return t.onTextChanged.Listen(callback)
}

func (t *TextBoxController) SelectionCount() int {
	return len(t.selections)
}

func (t *TextBoxController) Selection(index int) TextSelection {
	return t.selections[index]
}

func (t *TextBoxController) Selections() TextSelectionList {
	return append(TextSelectionList{}, t.selections...)
}

func (t *TextBoxController) SelectionText(index int) string {
	sel := t.selections[index]
	runes := t.text[sel.start:sel.end]
	return RuneArrayToString(runes)
}

func (t *TextBoxController) SelectionLineText(index int) string {
	sel := t.selections[index]
	line := t.LineIndex(sel.start)
	runes := t.text[t.LineStart(line):t.LineEnd(line)]
	return RuneArrayToString(runes)
}

func (t *TextBoxController) Caret(index int) int {
	return t.selections[index].Caret()
}

func (t *TextBoxController) Carets() []int {
	l := make([]int, len(t.selections))
	for i, s := range t.selections {
		l[i] = s.Caret()
	}
	return l
}

func (t *TextBoxController) FirstCaret() int {
	return t.Caret(0)
}

func (t *TextBoxController) LastCaret() int {
	return t.Caret(t.SelectionCount() - 1)
}

func (t *TextBoxController) FirstSelection() TextSelection {
	return t.Selection(0)
}

func (t *TextBoxController) LastSelection() TextSelection {
	return t.Selection(t.SelectionCount() - 1)
}

func (t *TextBoxController) LineCount() int {
	return len(t.lineStarts)
}

func (t *TextBoxController) Line(lineNo int) string {
	return RuneArrayToString(t.LineRunes(lineNo))
}

func (t *TextBoxController) LineRunes(lineNo int) []rune {
	start := t.LineStart(lineNo)
	end := t.LineEnd(lineNo)
	return t.text[start:end]
}

func (t *TextBoxController) LineStart(lineNo int) int {
	if t.LineCount() == 0 {
		return 0
	}
	return t.lineStarts[lineNo]
}

func (t *TextBoxController) LineEnd(lineNo int) int {
	if t.LineCount() == 0 {
		return 0
	}
	return t.lineEnds[lineNo]
}

func (t *TextBoxController) LineIndent(lineNo int) int {
	start, end := t.LineStart(lineNo), t.LineEnd(lineNo)
	result := end - start
	for i := 0; i < result; i++ {
		if !unicode.IsSpace(t.text[i+start]) {
			return i
		}
	}
	return result
}

func (t *TextBoxController) LineIndex(index int) int {
	return sort.Search(
		len(t.lineStarts),
		func(i int) bool {
			return index <= t.lineEnds[i]
		},
	)
}

func (t *TextBoxController) Text() string {
	return RuneArrayToString(t.text)
}

func (t *TextBoxController) TextRange(s, e int) string {
	return RuneArrayToString(t.text[s:e])
}

func (t *TextBoxController) TextRunes() []rune {
	return t.text
}

func (t *TextBoxController) SetText(text string) {
	t.SetTextRunes(StringToRuneArray(text))
}

func (t *TextBoxController) SetTextRunes(runes []rune) {
	t.setTextRunesNoEvent(runes)
	t.textEdited([]TextBoxEdit{})
}

func (t *TextBoxController) SetTextEdits(runes []rune, edits []TextBoxEdit) {
	t.setTextRunesNoEvent(runes)
	t.textEdited(edits)
}

func (t *TextBoxController) IndexFirst(index int) int {
	return 0
}

func (t *TextBoxController) IndexLast(index int) int {
	return len(t.text)
}

func (t *TextBoxController) IndexLeft(index int) int {
	return math.Max(index-1, 0)
}

func (t *TextBoxController) IndexRight(index int) int {
	return math.Min(index+1, len(t.text))
}

func (t *TextBoxController) IndexWordLeft(index int) int {
	index--
	if index >= 0 {
		wasInWord := t.RuneInWord(t.text[index])
		for index > 0 {
			isInWord := t.RuneInWord(t.text[index-1])
			if isInWord != wasInWord {
				return index
			}
			wasInWord = isInWord
			index--
		}
	}
	return 0
}

func (t *TextBoxController) IndexWordRight(index int) int {
	if index < len(t.text) {
		wasInWord := t.RuneInWord(t.text[index])
		for index < len(t.text)-1 {
			index++
			isInWord := t.RuneInWord(t.text[index])
			if isInWord != wasInWord {
				return index
			}
			wasInWord = isInWord
		}
	}
	return len(t.text)
}

func (t *TextBoxController) IndexUp(index int) int {
	l := t.LineIndex(index)
	x := index - t.LineStart(l)
	if l > 0 {
		return math.Min(t.LineStart(l-1)+x, t.LineEnd(l-1))
	} else {
		return 0
	}
}

func (t *TextBoxController) IndexDown(index int) int {
	l := t.LineIndex(index)
	x := index - t.LineStart(l)
	if l < t.LineCount()-1 {
		return math.Min(t.LineStart(l+1)+x, t.LineEnd(l+1))
	} else {
		return t.LineEnd(l)
	}
}

func (t *TextBoxController) IndexHome(index int) int {
	l := t.LineIndex(index)
	s := t.LineStart(l)
	x := index - s
	indent := t.LineIndent(l)
	if x > indent {
		return s + indent
	} else {
		return s
	}
}

func (t *TextBoxController) IndexEnd(index int) int {
	return t.LineEnd(t.LineIndex(index))
}

type SelectionTransform func(int) int

func (t *TextBoxController) ClearSelections() {
	t.storeCaretLocationsNextEdit = true
	t.SetCaret(t.Caret(0))
}

func (t *TextBoxController) SetCaret(index int) {
	t.storeCaretLocationsNextEdit = true
	t.selections = TextSelectionList{}
	t.AddCaret(index)
}

func (t *TextBoxController) AddCaret(index int) {
	t.storeCaretLocationsNextEdit = true
	t.AddSelection(TextSelection{index, index, false})
}

func (t *TextBoxController) AddSelection(selection TextSelection) {
	t.storeCaretLocationsNextEdit = true
	interval.Merge(&t.selections, selection)
	t.onSelectionChanged.Fire()
}

func (t *TextBoxController) SetSelection(selection TextSelection) {
	t.storeCaretLocationsNextEdit = true
	t.selections = []TextSelection{selection}
	t.onSelectionChanged.Fire()
}

func (t *TextBoxController) SetSelections(list TextSelectionList) {
	t.storeCaretLocationsNextEdit = true
	t.selections = list
	if len(list) == 0 {
		t.AddCaret(0)
	} else {
		t.onSelectionChanged.Fire()
	}
}

func (t *TextBoxController) SelectAll() {
	t.storeCaretLocationsNextEdit = true
	t.SetSelection(TextSelection{0, len(t.text), false})
}

func (t *TextBoxController) RestorePreviousSelections() {
	if t.locationHistoryIndex == len(t.locationHistory) {
		t.StoreCaretLocations()
		t.locationHistoryIndex--
	}

	if t.locationHistoryIndex > 0 {
		t.locationHistoryIndex--
		locations := t.locationHistory[t.locationHistoryIndex]
		t.selections = make(TextSelectionList, len(locations))
		for i, l := range locations {
			t.selections[i] = TextSelection{l, l, false}
		}
		t.onSelectionChanged.Fire()
	}
}

func (t *TextBoxController) RestoreNextSelections() {
	if t.locationHistoryIndex < len(t.locationHistory)-1 {
		t.locationHistoryIndex++
		locations := t.locationHistory[t.locationHistoryIndex]
		t.selections = make(TextSelectionList, len(locations))
		for i, l := range locations {
			t.selections[i] = TextSelection{l, l, false}
		}
		t.onSelectionChanged.Fire()
	}
}

func (t *TextBoxController) AddCarets(transform SelectionTransform) {
	t.storeCaretLocationsNextEdit = true
	up := t.selections.Transform(0, transform)
	for _, s := range up {
		interval.Merge(&t.selections, s)
	}
	t.onSelectionChanged.Fire()
}

func (t *TextBoxController) GrowSelections(transform SelectionTransform) {
	t.storeCaretLocationsNextEdit = true
	t.selections = t.selections.TransformCarets(0, transform)
	t.onSelectionChanged.Fire()
}

func (t *TextBoxController) MoveSelections(transform SelectionTransform) {
	t.storeCaretLocationsNextEdit = true
	t.selections = t.selections.Transform(0, transform)
	t.onSelectionChanged.Fire()
}

func (t *TextBoxController) AddCaretsUp()       { t.AddCarets(t.IndexUp) }
func (t *TextBoxController) AddCaretsDown()     { t.AddCarets(t.IndexDown) }
func (t *TextBoxController) SelectFirst()       { t.GrowSelections(t.IndexFirst) }
func (t *TextBoxController) SelectLast()        { t.GrowSelections(t.IndexLast) }
func (t *TextBoxController) SelectLeft()        { t.GrowSelections(t.IndexLeft) }
func (t *TextBoxController) SelectRight()       { t.GrowSelections(t.IndexRight) }
func (t *TextBoxController) SelectUp()          { t.GrowSelections(t.IndexUp) }
func (t *TextBoxController) SelectDown()        { t.GrowSelections(t.IndexDown) }
func (t *TextBoxController) SelectHome()        { t.GrowSelections(t.IndexHome) }
func (t *TextBoxController) SelectEnd()         { t.GrowSelections(t.IndexEnd) }
func (t *TextBoxController) SelectLeftByWord()  { t.GrowSelections(t.IndexWordLeft) }
func (t *TextBoxController) SelectRightByWord() { t.GrowSelections(t.IndexWordRight) }
func (t *TextBoxController) MoveFirst()         { t.MoveSelections(t.IndexFirst) }
func (t *TextBoxController) MoveLast()          { t.MoveSelections(t.IndexLast) }
func (t *TextBoxController) MoveLeft()          { t.MoveSelections(t.IndexLeft) }
func (t *TextBoxController) MoveRight()         { t.MoveSelections(t.IndexRight) }
func (t *TextBoxController) MoveUp()            { t.MoveSelections(t.IndexUp) }
func (t *TextBoxController) MoveDown()          { t.MoveSelections(t.IndexDown) }
func (t *TextBoxController) MoveLeftByWord()    { t.MoveSelections(t.IndexWordLeft) }
func (t *TextBoxController) MoveRightByWord()   { t.MoveSelections(t.IndexWordRight) }
func (t *TextBoxController) MoveHome()          { t.MoveSelections(t.IndexHome) }
func (t *TextBoxController) MoveEnd()           { t.MoveSelections(t.IndexEnd) }

func (t *TextBoxController) Delete() {
	t.maybeStoreCaretLocations()
	text := t.text
	var edits []TextBoxEdit

	for index := len(t.selections) - 1; index >= 0; index-- {
		selection := t.selections[index]
		if selection.start == selection.end && selection.end < len(t.text) {
			copy(text[selection.start:], text[selection.start+1:])
			text = text[:len(text)-1]
			edits = append(edits, TextBoxEdit{selection.start, -1})
		} else {
			copy(text[selection.start:], text[selection.end:])
			l := selection.Length()
			text = text[:len(text)-l]
			edits = append(edits, TextBoxEdit{selection.start, -l})
		}
		t.selections[index] = TextSelection{selection.end, selection.end, false}
	}

	t.SetTextEdits(text, edits)
}

func (t *TextBoxController) Backspace() {
	t.maybeStoreCaretLocations()
	text := t.text
	var edits []TextBoxEdit

	for index := len(t.selections) - 1; index >= 0; index-- {
		selection := t.selections[index]
		if selection.start == selection.end && selection.start > 0 {
			copy(text[selection.start-1:], text[selection.start:])
			text = text[:len(text)-1]
			edits = append(edits, TextBoxEdit{selection.start - 1, -1})
		} else {
			copy(text[selection.start:], text[selection.end:])
			l := selection.Length()
			text = text[:len(text)-l]
			edits = append(edits, TextBoxEdit{selection.start - 1, -l})
		}
		t.selections[index] = TextSelection{selection.end, selection.end, false}
	}

	t.SetTextEdits(text, edits)
}

func (t *TextBoxController) ReplaceAll(text string) {
	t.Replace(func(TextSelection) string { return text })
}

func (t *TextBoxController) ReplaceAllRunes(runes []rune) {
	t.ReplaceRunes(func(TextSelection) []rune { return runes })
}

func (t *TextBoxController) Replace(callback func(selection TextSelection) string) {
	t.ReplaceRunes(func(selection TextSelection) []rune { return StringToRuneArray(callback(selection)) })
}

func (t *TextBoxController) ReplaceRunes(callback func(selection TextSelection) []rune) {
	t.maybeStoreCaretLocations()
	text := t.text
	edit := TextBoxEdit{}
	var edits []TextBoxEdit

	for index := len(t.selections) - 1; index >= 0; index-- {
		selection := t.selections[index]
		text, edit = t.ReplaceAt(text, selection.start, selection.end, callback(selection))
		edits = append(edits, edit)
	}

	t.setTextRunesNoEvent(text)
	t.textEdited(edits)
}

func (t *TextBoxController) ReplaceAt(text []rune, start, end int, replacement []rune) ([]rune, TextBoxEdit) {
	replacementLen := len(replacement)
	delta := replacementLen - (end - start)
	if delta > 0 {
		text = append(text, make([]rune, delta)...)
	}

	copy(text[end+delta:], text[end:])
	copy(text[start:], replacement)

	if delta < 0 {
		text = text[:len(text)+delta]
	}

	return text, TextBoxEdit{start, delta}
}

func (t *TextBoxController) ReplaceWithNewline() {
	t.ReplaceAll("\n")
	t.Deselect(false)
}

func (t *TextBoxController) ReplaceWithNewlineKeepIndent() {
	t.Replace(
		func(selection TextSelection) string {
			start, _ := selection.Range()
			indent := t.LineIndent(t.LineIndex(start))
			return "\n" + strings.Repeat(" ", indent)
		},
	)
	t.Deselect(false)
}

func (t *TextBoxController) IndentSelection(tabWidth int) {
	tab := make([]rune, tabWidth)
	for i := range tab {
		tab[i] = ' '
	}

	text := t.text
	edit := TextBoxEdit{}
	var edits []TextBoxEdit

	lastLine := -1
	for index := len(t.selections) - 1; index >= 0; index-- {
		selection := t.selections[index]

		lineStart, lineEnd := t.LineIndex(selection.start), t.LineIndex(selection.end)
		if lastLine == lineEnd {
			lineEnd--
		}

		for l := lineEnd; l >= lineStart; l-- {
			ls := t.LineStart(l)
			text, edit = t.ReplaceAt(text, ls, ls, tab)
			edits = append(edits, edit)
		}

		lastLine = lineStart
	}

	t.SetTextEdits(text, edits)
}

func (t *TextBoxController) UnindentSelection(tabWidth int) {
	text := t.text
	edit := TextBoxEdit{}
	var edits []TextBoxEdit

	lastLine := -1
	for index := len(t.selections) - 1; index >= 0; index-- {
		selection := t.selections[index]
		lineStart, lineEnd := t.LineIndex(selection.start), t.LineIndex(selection.end)
		if lastLine == lineEnd {
			lineEnd--
		}

		for l := lineEnd; l >= lineStart; l-- {
			c := math.Min(t.LineIndent(l), tabWidth)
			if c > 0 {
				ls := t.LineStart(l)
				text, edit = t.ReplaceAt(text, ls, ls+c, []rune{})
				edits = append(edits, edit)
			}
		}
		lastLine = lineStart
	}
	t.SetTextEdits(text, edits)
}

func (t *TextBoxController) RuneInWord(theRune rune) bool {
	switch {
	case unicode.IsLetter(theRune), unicode.IsNumber(theRune), theRune == '_':
		return true
	default:
		return false
	}
}

func (t *TextBoxController) WordAt(index int) (int, int) {
	text := t.text
	start, end := index, index
	for start > 0 && t.RuneInWord(text[start-1]) {
		start--
	}
	for end < len(t.text) && t.RuneInWord(text[end]) {
		end++
	}
	return start, end
}

func (t *TextBoxController) Deselect(moveCaretToStart bool) bool {
	deselected := false
	for index, selection := range t.selections {
		if selection.start == selection.end {
			continue
		}
		deselected = true
		if moveCaretToStart {
			selection.end = selection.start
		} else {
			selection.start = selection.end
		}
		t.selections[index] = selection
	}
	if deselected {
		t.onSelectionChanged.Fire()
	}
	return deselected
}

func (t *TextBoxController) LineAndRow(index int) (int, int) {
	line := t.LineIndex(index)
	row := index - t.LineStart(line)
	return line, row
}
