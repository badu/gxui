// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mixins

import (
	"github.com/badu/gxui"
)

type AttachablePart struct {
	onAttach gxui.Event
	onDetach gxui.Event
	attached bool
}

func (a *AttachablePart) Init() {}

func (a *AttachablePart) Attached() bool {
	return a.attached
}

func (a *AttachablePart) Attach() {
	if a.attached {
		panic("Control already attached")
	}
	a.attached = true
	if a.onAttach != nil {
		a.onAttach.Fire()
	}
}

func (a *AttachablePart) Detach() {
	if !a.attached {
		panic("Control already detached")
	}
	a.attached = false
	if a.onDetach != nil {
		a.onDetach.Fire()
	}
}

func (a *AttachablePart) OnAttach(callback func()) gxui.EventSubscription {
	if a.onAttach == nil {
		a.onAttach = gxui.CreateEvent(func() {})
	}
	return a.onAttach.Listen(callback)
}

func (a *AttachablePart) OnDetach(callback func()) gxui.EventSubscription {
	if a.onDetach == nil {
		a.onDetach = gxui.CreateEvent(func() {})
	}
	return a.onDetach.Listen(callback)
}
