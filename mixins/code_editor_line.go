// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/interval"
	"github.com/badu/gxui/math"
)

type CodeEditorLineOuter interface {
	gxui.DefaultTextBoxLineOuter
	PaintBackgroundSpans(c gxui.Canvas, info CodeEditorLinePaintInfo)
	PaintGlyphs(c gxui.Canvas, info CodeEditorLinePaintInfo)
	PaintBorders(c gxui.Canvas, info CodeEditorLinePaintInfo)
}

type CodeEditorLinePaintInfo struct {
	LineSpan     interval.IntData
	Runes        []rune
	GlyphOffsets []math.Point
	GlyphWidth   int
	LineHeight   int
	Font         gxui.Font
}

// CodeEditorLine
type CodeEditorLine struct {
	DefaultTextBoxLine
	outer  CodeEditorLineOuter
	editor *CodeEditor
}

func (l *CodeEditorLine) Init(outer CodeEditorLineOuter, theme gxui.Theme, ce *CodeEditor, lineIndex int) {
	l.DefaultTextBoxLine.Init(outer, theme, &ce.TextBox, lineIndex)
	l.outer = outer
	l.editor = ce
}

func (t *CodeEditorLine) PaintBackgroundSpans(canvas gxui.Canvas, info CodeEditorLinePaintInfo) {
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
					canvas.DrawRoundedRect(r, 3, 3, 3, 3, gxui.TransparentPen, gxui.Brush{Color: color})
				})
				interval.Remove(&remaining, span)
			}
		}
	}
}

func (t *CodeEditorLine) PaintGlyphs(canvas gxui.Canvas, info CodeEditorLinePaintInfo) {
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
		s, e := span.Span()
		s, e = s-start, e-start
		canvas.DrawRunes(font, runes[s:e], offsets[s:e], t.editor.textColor)
	}
}

func (t *CodeEditorLine) PaintBorders(canvas gxui.Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	for _, layer := range t.editor.layers {
		if layer != nil && layer.BorderColor() != nil {
			color := *layer.BorderColor()
			interval.Visit(layer.Spans(), info.LineSpan, func(vs, ve uint64, _ int) {
				s, e := vs-start, ve-start
				r := math.CreateRect(offsets[s].X, 0, offsets[e-1].X+info.GlyphWidth, info.LineHeight)
				canvas.DrawRoundedRect(r, 3, 3, 3, 3, gxui.CreatePen(0.5, color), gxui.TransparentBrush)
			})
		}
	}
}

// DefaultTextBoxLine overrides
func (t *CodeEditorLine) Paint(canvas gxui.Canvas) {
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
			&gxui.TextBlock{
				Runes:     runes,
				AlignRect: rect,
				H:         gxui.AlignLeft,
				V:         gxui.AlignMiddle,
			},
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
		t.outer.PaintBackgroundSpans(canvas, info)

		// Selections
		if t.textbox.HasFocus() {
			t.outer.PaintSelections(canvas)
		}

		// Glyphs
		t.outer.PaintGlyphs(canvas, info)

		// Borders
		t.outer.PaintBorders(canvas, info)
	}

	// Carets
	if t.textbox.HasFocus() {
		t.outer.PaintCarets(canvas)
	}
}
