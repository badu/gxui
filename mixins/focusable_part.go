// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type FocusableOuter interface{}

type FocusablePart struct {
	outer         FocusableOuter
	focusable     bool
	hasFocus      bool
	onGainedFocus gxui.Event
	onLostFocus   gxui.Event
}

func (f *FocusablePart) Init(outer FocusableOuter) {
	f.outer = outer
	f.focusable = true
}

// gxui.Control compliance
func (f *FocusablePart) IsFocusable() bool {
	return f.focusable
}

func (f *FocusablePart) HasFocus() bool {
	return f.hasFocus
}

func (f *FocusablePart) SetFocusable(bool) {
	f.focusable = true
}

func (f *FocusablePart) OnGainedFocus(callback func()) gxui.EventSubscription {
	if f.onGainedFocus == nil {
		f.onGainedFocus = gxui.CreateEvent(f.GainedFocus)
	}
	return f.onGainedFocus.Listen(callback)
}

func (f *FocusablePart) OnLostFocus(callback func()) gxui.EventSubscription {
	if f.onLostFocus == nil {
		f.onLostFocus = gxui.CreateEvent(f.LostFocus)
	}
	return f.onLostFocus.Listen(callback)
}

func (f *FocusablePart) GainedFocus() {
	f.hasFocus = true
	if f.onGainedFocus != nil {
		f.onGainedFocus.Fire()
	}
}

func (f *FocusablePart) LostFocus() {
	f.hasFocus = false
	if f.onLostFocus != nil {
		f.onLostFocus.Fire()
	}
}
