// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type Label interface {
	Control
	Text() string
	SetText(text string)
	Font() Font
	SetFont(font Font)
	Color() Color
	SetColor(color Color)
	Multiline() bool
	SetMultiline(bool)
	SetHorizontalAlignment(HorizontalAlignment)
	HorizontalAlignment() HorizontalAlignment
	SetVerticalAlignment(VerticalAlignment)
	VerticalAlignment() VerticalAlignment
}
