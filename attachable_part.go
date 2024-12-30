// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type AttachablePart struct {
	onAttach Event
	onDetach Event
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

func (a *AttachablePart) OnAttach(callback func()) EventSubscription {
	if a.onAttach == nil {
		a.onAttach = CreateEvent(func() {})
	}
	return a.onAttach.Listen(callback)
}

func (a *AttachablePart) OnDetach(callback func()) EventSubscription {
	if a.onDetach == nil {
		a.onDetach = CreateEvent(func() {})
	}
	return a.onDetach.Listen(callback)
}
