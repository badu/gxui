// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/interval"
	"github.com/badu/gxui/pkg/math"
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
	DefaultTextBoxLine
	parent CodeEditorLineParent
	editor *CodeEditor
}

func (l *CodeEditorLine) Init(parent CodeEditorLineParent, editor *CodeEditor, lineIndex int) {
	l.DefaultTextBoxLine.Init(parent, &editor.TextBox, lineIndex)
	l.parent = parent
	l.editor = editor
}

func (l *CodeEditorLine) PaintBackgroundSpans(canvas Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	remaining := interval.IntDataList{info.LineSpan}
	for _, layer := range l.editor.layers {
		if layer != nil && layer.BackgroundColor() != nil {
			color := *layer.BackgroundColor()
			for _, span := range layer.Spans().Overlaps(info.LineSpan) {
				interval.Visit(
					&remaining,
					span,
					func(vs, ve uint64, _ int) {
						s, e := vs-start, ve-start
						r := math.CreateRect(offsets[s].X, 0, offsets[e-1].X+info.GlyphWidth, info.LineHeight)
						canvas.DrawRoundedRect(r, 3, 3, 3, 3, TransparentPen, Brush{Color: color})
					},
				)
				interval.Remove(&remaining, span)
			}
		}
	}
}

func (l *CodeEditorLine) PaintGlyphs(canvas Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	runes, offsets, font := info.Runes, info.GlyphOffsets, info.Font
	remaining := interval.IntDataList{info.LineSpan}
	for _, layer := range l.editor.layers {
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
		canvas.DrawRunes(font, runes[spanStart:spanEnd], offsets[spanStart:spanEnd], l.editor.textColor)
	}
}

func (l *CodeEditorLine) PaintBorders(canvas Canvas, info CodeEditorLinePaintInfo) {
	start, _ := info.LineSpan.Span()
	offsets := info.GlyphOffsets
	for _, layer := range l.editor.layers {
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
func (l *CodeEditorLine) Paint(canvas Canvas) {
	font := l.editor.font
	rect := l.Size().Rect().OffsetX(l.caretWidth)
	controller := l.editor.controller
	runes := controller.LineRunes(l.lineIndex)
	start := controller.LineStart(l.lineIndex)
	end := controller.LineEnd(l.lineIndex)

	if start != end {
		lineSpan := interval.CreateIntData(start, end, nil)

		lineHeight := l.Size().Height
		glyphWidth := font.GlyphMaxSize().Width
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
		l.parent.PaintBackgroundSpans(canvas, info)

		// Selections
		if l.textbox.HasFocus() {
			l.parent.PaintSelections(canvas)
		}

		// Glyphs
		l.parent.PaintGlyphs(canvas, info)

		// Borders
		l.parent.PaintBorders(canvas, info)
	}

	// Carets
	if l.textbox.HasFocus() {
		l.parent.PaintCarets(canvas)
	}
}

func (l *CodeEditorLine) ContainsPoint(point math.Point) bool {
	return l.IsVisible() && l.Size().Rect().Contains(point)
}
