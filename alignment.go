// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type HAlign int

const (
	AlignLeft HAlign = iota
	AlignCenter
	AlignRight
)

func (a HAlign) AlignLeft() bool   { return a == AlignLeft }
func (a HAlign) AlignCenter() bool { return a == AlignCenter }
func (a HAlign) AlignRight() bool  { return a == AlignRight }

type VAlign int

const (
	AlignTop VAlign = iota
	AlignMiddle
	AlignBottom
)

func (a VAlign) AlignTop() bool    { return a == AlignTop }
func (a VAlign) AlignMiddle() bool { return a == AlignMiddle }
func (a VAlign) AlignBottom() bool { return a == AlignBottom }
