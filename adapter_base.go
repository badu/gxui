// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type AdapterBase struct {
	onDataChanged  Event
	onDataReplaced Event
}

func (a *AdapterBase) DataChanged(recreateControls bool) {
	if a.onDataChanged == nil {
		return
	}

	a.onDataChanged.Emit(recreateControls)
}

func (a *AdapterBase) DataReplaced() {
	if a.onDataReplaced == nil {
		return
	}

	a.onDataReplaced.Emit()
}

func (a *AdapterBase) OnDataChanged(callback func(recreateControls bool)) EventSubscription {
	if a.onDataChanged == nil {
		a.onDataChanged = CreateEvent(func(recreateControls bool) {})
	}

	return a.onDataChanged.Listen(callback)
}

func (a *AdapterBase) OnDataReplaced(callback func()) EventSubscription {
	if a.onDataReplaced == nil {
		a.onDataReplaced = CreateEvent(func() {})
	}

	return a.onDataReplaced.Listen(callback)
}
