// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/pkg/interval"
)

type TextSelectionList []TextSelection

func (l TextSelectionList) Transform(from int, transform func(index int) int) TextSelectionList {
	result := TextSelectionList{}
	for _, item := range l {
		start := item.start
		end := item.end
		if start >= from {
			start = transform(start)
		}
		if end >= from {
			end = transform(end)
		}
		interval.Merge(&result, TextSelection{start, end, item.caretAtStart})
	}
	return result
}

func (l TextSelectionList) TransformCarets(from int, transform func(index int) int) TextSelectionList {
	result := TextSelectionList{}
	for _, item := range l {
		if item.caretAtStart && item.start >= from {
			item.start = transform(item.start)
		} else if item.end >= from {
			item.end = transform(item.end)
		}
		if item.start > item.end {
			tmp := item.start
			item.start = item.end
			item.end = tmp
			item.caretAtStart = !item.caretAtStart
		}
		interval.Merge(&result, item)
	}
	return result
}

func (l TextSelectionList) Len() int {
	return len(l)
}

func (l TextSelectionList) Cap() int {
	return cap(l)
}

func (l *TextSelectionList) SetLen(len int) {
	*l = (*l)[:len]
}

func (l *TextSelectionList) GrowTo(length, capacity int) {
	old := *l
	*l = make(TextSelectionList, length, capacity)
	copy(*l, old)
}

func (l TextSelectionList) Copy(to, from, count int) {
	copy(l[to:to+count], l[from:from+count])
}

func (l TextSelectionList) GetInterval(index int) (start, end uint64) {
	return l[index].Span()
}

func (l TextSelectionList) SetInterval(index int, start, end uint64) {
	l[index].start = int(start)
	l[index].end = int(end)
}

func (l TextSelectionList) MergeData(index int, i interval.Node) {
	l[index].caretAtStart = i.(TextSelection).caretAtStart
}
