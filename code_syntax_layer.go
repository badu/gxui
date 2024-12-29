// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/interval"
	"github.com/badu/gxui/math"
)

type CodeSyntaxLayer struct {
	spans           interval.IntDataList
	color           *Color
	backgroundColor *Color
	borderColor     *Color
	data            interface{}
}

func CreateCodeSyntaxLayer() *CodeSyntaxLayer { return &CodeSyntaxLayer{} }

func (l *CodeSyntaxLayer) Clear() {
	l.spans = interval.IntDataList{}
}

func (l *CodeSyntaxLayer) UpdateSpans(runeCount int, edits []TextBoxEdit) {
	pMin := 0
	pMax := runeCount
	for _, edit := range edits {
		if l == nil { // TODO : @Badu - why?
			continue
		}

		for index, span := range l.spans {
			at := edit.At
			start, end := span.Range()
			if start >= at {
				start = math.Clamp(start+edit.Delta, pMin, pMax)
			}

			if end > at {
				end = math.Clamp(end+edit.Delta, pMin, pMax)
			}

			if end < start {
				end = start
			}

			l.spans[index] = interval.CreateIntData(start, end, span.Data())
		}
	}
}

func (l *CodeSyntaxLayer) Add(start, count int) {
	l.AddData(start, count, nil)
}

func (l *CodeSyntaxLayer) AddData(start, count int, data interface{}) {
	span := interval.CreateIntData(start, start+count, data)
	interval.Replace(&l.spans, span)
}

func (l *CodeSyntaxLayer) AddSpan(span interval.IntData) {
	interval.Replace(&l.spans, span)
}

func (l *CodeSyntaxLayer) Spans() interval.IntDataList {
	return l.spans
}

func (l *CodeSyntaxLayer) SpanAt(runeIndex int) *interval.IntData {
	idx := interval.IndexOf(&l.spans, uint64(runeIndex))
	if idx >= 0 {
		return &l.spans[idx]
	} else {
		return nil
	}
}

func (l *CodeSyntaxLayer) Color() *Color {
	return l.color
}

func (l *CodeSyntaxLayer) ClearColor() {
	l.color = nil
}

func (l *CodeSyntaxLayer) SetColor(color Color) {
	l.color = &color
}

func (l *CodeSyntaxLayer) BackgroundColor() *Color {
	return l.backgroundColor
}

func (l *CodeSyntaxLayer) ClearBackgroundColor() {
	l.backgroundColor = nil
}

func (l *CodeSyntaxLayer) SetBackgroundColor(color Color) {
	l.backgroundColor = &color
}

func (l *CodeSyntaxLayer) BorderColor() *Color {
	return l.borderColor
}

func (l *CodeSyntaxLayer) ClearBorderColor() {
	l.borderColor = nil
}

func (l *CodeSyntaxLayer) SetBorderColor(color Color) {
	l.borderColor = &color
}

func (l *CodeSyntaxLayer) Data() interface{} {
	return l.data
}

func (l *CodeSyntaxLayer) SetData(data interface{}) {
	l.data = data
}

type CodeSyntaxLayers []*CodeSyntaxLayer

func (l *CodeSyntaxLayers) Get(idx int) *CodeSyntaxLayer {
	if len(*l) <= idx {
		old := *l
		*l = make(CodeSyntaxLayers, idx+1)
		copy(*l, old)
	}
	layer := (*l)[idx]
	if layer == nil {
		layer = &CodeSyntaxLayer{}
		(*l)[idx] = layer
	}
	return layer
}

func (l *CodeSyntaxLayers) Clear() {
	*l = CodeSyntaxLayers{}
}
