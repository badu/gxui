// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package interval

import "sort"

type Node interface {
	Span() (start, end uint64)
}

type RList interface {
	Len() int
	GetInterval(index int) (start, end uint64)
	SetInterval(index int, start, end uint64)
}

type List interface {
	Len() int
	GetInterval(index int) (start, end uint64)
	SetInterval(index int, start, end uint64)
	Copy(to, from, count int)
	Cap() int
	SetLen(len int)
	GrowTo(length, capacity int)
}

type ExtendedList interface {
	MergeData(index int, i Node)
}

type intersection struct {
	overlap        int
	lowIndex       int
	lowStart       uint64
	lowEnd         uint64
	intersectsLow  bool
	highIndex      int
	highStart      uint64
	highEnd        uint64
	intersectsHigh bool
}

const (
	minSpace = 3 // Max growth of 2 plus one slot for temporary
	minCap   = 5
)

func Merge(target List, node Node) {
	start, end := node.Span()
	s := intersection{}
	s.intersect(target, start, end)
	adjust(target, s.lowIndex, 1-s.overlap)
	if s.intersectsLow {
		start = s.lowStart
	}
	if s.intersectsHigh {
		end = s.highEnd
	}
	target.SetInterval(s.lowIndex, start, end)
	if dl, ok := target.(ExtendedList); ok {
		dl.MergeData(s.lowIndex, node)
	}
}

func Replace(target List, node Node) {
	start, end := node.Span()
	index, start, end := replace(target, start, end, true)
	target.SetInterval(index, start, end)
	if dl, ok := target.(ExtendedList); ok {
		dl.MergeData(index, node)
	}
}

func Remove(target List, node Node) {
	start, end := node.Span()
	replace(target, start, end, false)
}

func Intersect(target RList, node Node) (first, count int) {
	start, end := node.Span()
	s := intersection{}
	s.intersect(target, start, end)
	return s.lowIndex, s.overlap
}

type Visitor func(start, end uint64, index int)

func Visit(target RList, node Node, visitorFn Visitor) {
	start, end := node.Span()
	s := intersection{}
	s.intersect(target, start, end)
	for index := s.lowIndex; index < s.lowIndex+s.overlap; index++ {
		s, e := target.GetInterval(index)
		if s < start {
			s = start
		}
		if e > end {
			e = end
		}
		visitorFn(s, e, index)
	}
}

func Contains(target RList, index uint64) bool {
	return IndexOf(target, index) >= 0
}

func IndexOf(target RList, p uint64) int {
	index := sort.Search(
		target.Len(),
		func(at int) bool {
			iStart, _ := target.GetInterval(at)
			return p < iStart
		},
	)
	index--
	if index >= 0 {
		_, iEnd := target.GetInterval(index)
		if p < iEnd {
			return index
		}
	}
	return -1
}

func FindStart(target RList, at int, start uint64) bool {
	_, end := target.GetInterval(at)
	return start < end
}

func FindEnd(target RList, at int, end uint64) bool {
	start, _ := target.GetInterval(at)
	return end <= start
}

type Searcher func(target RList, at int, v uint64) bool

func Search(target RList, v uint64, searcherFn Searcher) int {
	i, j := 0, target.Len()
	for i < j {
		h := i + (j-i)/2
		if !searcherFn(target, h, v) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

func (s *intersection) intersect(target RList, start, end uint64) {
	beforeLen := Search(target, start, FindStart)
	afterIndex := Search(target, end, FindEnd)
	if afterIndex < beforeLen {
		afterIndex, beforeLen = beforeLen, afterIndex
	}
	s.lowIndex = beforeLen
	s.highIndex = afterIndex - 1
	s.overlap = afterIndex - beforeLen
	s.intersectsLow = false
	s.intersectsHigh = false
	if s.overlap > 0 {
		s.lowStart, s.lowEnd = target.GetInterval(s.lowIndex)
		s.intersectsLow = s.lowStart < start
		s.highStart, s.highEnd = target.GetInterval(s.highIndex)
		s.intersectsHigh = end < s.highEnd
	}
}

func adjust(target List, at, delta int) {
	if delta == 0 {
		return
	}
	oldLen := target.Len()
	newLen := oldLen + delta
	if delta > 0 {
		cap := target.Cap()
		if cap < newLen {
			newCap := newLen + (newLen >> 1)
			target.GrowTo(newLen, newCap)
		} else {
			target.SetLen(newLen)
		}
	}
	copyStart := at - delta
	copyTo := at
	if copyStart < 0 {
		copyTo -= copyStart
		copyStart = 0
	}
	target.Copy(copyTo, copyStart, newLen-copyTo)
	if delta < 0 {
		target.SetLen(newLen)
	}
}

func replace(target List, start, end uint64, add bool) (int, uint64, uint64) {
	s := intersection{}
	s.intersect(target, start, end)
	if s.overlap == 0 {
		if add {
			adjust(target, s.lowIndex, 1)
		}
		return s.lowIndex, start, end
	}

	insertLen := 0
	insertPoint := s.lowIndex
	if s.intersectsLow {
		s.lowEnd = start
		insertLen++
		insertPoint++
	}
	if add {
		insertLen++
	}
	if s.intersectsHigh {
		s.highStart = end
		insertLen++
	}
	delta := insertLen - s.overlap
	adjust(target, insertPoint, delta)
	if s.intersectsLow {
		target.SetInterval(s.lowIndex, s.lowStart, s.lowEnd)
	}
	if s.intersectsHigh {
		target.SetInterval(s.lowIndex+insertLen-1, s.highStart, s.highEnd)
	}
	return insertPoint, start, end
}
