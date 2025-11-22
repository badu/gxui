// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"reflect"
	"sync"
)

type ChanneledEvent struct {
	channel chan func()
	base    EventBase
	sync.RWMutex
}

func CreateChanneledEvent(signature interface{}, channel chan func()) Event {
	result := &ChanneledEvent{channel: channel}
	result.base.init(signature)
	baseUnlisten := result.base.unlisten
	result.base.unlisten = func(id int) {
		result.RLock()
		baseUnlisten(id)
		result.RUnlock()
	}
	return result
}

func (e *ChanneledEvent) Emit(args ...interface{}) {
	e.base.VerifyArguments(args)
	e.channel <- func() {
		e.RLock()
		e.base.InvokeListeners(args)
		e.RUnlock()
	}
}

func (e *ChanneledEvent) Listen(listener interface{}) EventSubscription {
	e.Lock()
	res := e.base.Listen(listener)
	e.Unlock()
	return res
}

func (e *ChanneledEvent) ParameterTypes() []reflect.Type {
	return e.base.ParameterTypes()
}
