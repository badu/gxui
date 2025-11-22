// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/interval"
	"github.com/badu/gxui/math"
)

type CodeEditorLineParent interface {
	DefaultTextBoxLineParent
	PaintBackgroundSpans(c Canvas, info CodeEditorLinePaintInfo)
	PaintGlyphs(c Canvas, info CodeEditorLinePaintInfo)
	PaintBorders(c Canvas, info CodeEditorLinePaintInfo)
}

type CodeEditorLinePaintInfo struct {
	LineSpan     interval.IntData
	Font         Font
	Runes        []rune
	GlyphOffsets []math.Point
	GlyphWidth   int
	LineHeight   int
}

// CodeEditorLine
type CodeEditorLine struct {
	parent CodeEditorLineParent
	editor *CodeEditor
	DefaultTextBoxLine
}

func (l *CodeEditorLine) Init(parent CodeEditorLineParent, editor *CodeEditor, lineIndex int) {
	l.DefaultTextBoxLine.Init(parent, &editor.TextBox, lineIndex)
	l.parent = parent
	l.editor = editor
}

func (t *CodeEditorLine) PaintBackgroundSpans(canvas Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	remaining := interval.IntDataList{info.LineSpan}
	for _, layer := range t.editor.layers {
		if layer != nil && layer.BackgroundColor() != nil {
			color := *layer.BackgroundColor()
			for _, span := range layer.Spans().Overlaps(info.LineSpan) {
				interval.Visit(&remaining, span, func(vs, ve uint64, _ int) {
					s, e := vs-start, ve-start
					r := math.CreateRect(offsets[s].X, 0, offsets[e-1].X+info.GlyphWidth, info.LineHeight)
					canvas.DrawRoundedRect(r, 3, 3, 3, 3, TransparentPen, Brush{Color: color})
				})
				interval.Remove(&remaining, span)
			}
		}
	}
}

func (t *CodeEditorLine) PaintGlyphs(canvas Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	runes, offsets, font := info.Runes, info.GlyphOffsets, info.Font
	remaining := interval.IntDataList{info.LineSpan}
	for _, layer := range t.editor.layers {
		if layer != nil && layer.Color() != nil {
			color := *layer.Color()
			for _, span := range layer.Spans().Overlaps(info.LineSpan) {
				interval.Visit(&remaining, span, func(vs, ve uint64, _ int) {
					s, e := vs-start, ve-start
					canvas.DrawRunes(font, runes[s:e], offsets[s:e], color)
				})
				interval.Remove(&remaining, span)
			}
		}
	}

	for _, span := range remaining {
		spanStart, spanEnd := span.Span()
		spanStart, spanEnd = spanStart-start, spanEnd-start
		canvas.DrawRunes(font, runes[spanStart:spanEnd], offsets[spanStart:spanEnd], t.editor.textColor)
	}
}

func (t *CodeEditorLine) PaintBorders(canvas Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	for _, layer := range t.editor.layers {
		if layer != nil && layer.BorderColor() != nil {
			color := *layer.BorderColor()
			interval.Visit(layer.Spans(), info.LineSpan, func(vs, ve uint64, _ int) {
				s, e := vs-start, ve-start
				r := math.CreateRect(offsets[s].X, 0, offsets[e-1].X+info.GlyphWidth, info.LineHeight)
				canvas.DrawRoundedRect(r, 3, 3, 3, 3, CreatePen(0.5, color), TransparentBrush)
			})
		}
	}
}

// DefaultTextBoxLine overrides
func (t *CodeEditorLine) Paint(canvas Canvas) {
	font := t.editor.font
	rect := t.Size().Rect().OffsetX(t.caretWidth)
	controller := t.editor.controller
	runes := controller.LineRunes(t.lineIndex)
	start := controller.LineStart(t.lineIndex)
	end := controller.LineEnd(t.lineIndex)

	if start != end {
		lineSpan := interval.CreateIntData(start, end, nil)

		lineHeight := t.Size().H
		glyphWidth := font.GlyphMaxSize().W
		offsets := font.Layout(
			&TextBlock{Runes: runes, AlignRect: rect, H: AlignLeft, V: AlignMiddle},
		)

		info := CodeEditorLinePaintInfo{
			LineSpan:     lineSpan,
			Runes:        runes, // TODO gxui.TextBlock?
			GlyphOffsets: offsets,
			GlyphWidth:   glyphWidth,
			LineHeight:   lineHeight,
			Font:         font,
		}

		// Background
		t.parent.PaintBackgroundSpans(canvas, info)

		// Selections
		if t.textbox.HasFocus() {
			t.parent.PaintSelections(canvas)
		}

		// Glyphs
		t.parent.PaintGlyphs(canvas, info)

		// Borders
		t.parent.PaintBorders(canvas, info)
	}

	// Carets
	if t.textbox.HasFocus() {
		t.parent.PaintCarets(canvas)
	}
}
