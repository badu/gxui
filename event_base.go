// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"reflect"
)

type EventSubscription interface {
	Forget()
}

type Event interface {
	Emit(args ...interface{})
	Listen(event interface{}) EventSubscription
	ParameterTypes() []reflect.Type
}

type SimpleEvent struct {
	EventBase
}

func CreateEvent(signature interface{}) Event {
	result := &SimpleEvent{}
	result.init(signature)
	return result
}

type EventListener struct {
	Id       int
	Function reflect.Value
}

type eventBaseSubscription struct {
	event *EventBase
	id    int
}

func (s *eventBaseSubscription) Forget() {
	if s.event != nil {
		s.event.unlisten(s.id)
		s.event = nil
	}
}

type EventBase struct {
	unlisten   func(id int)
	paramTypes []reflect.Type
	isVariadic bool
	listeners  []EventListener
	nextId     int
}

func (e *EventBase) init(signature interface{}) {
	e.unlisten = func(id int) {
		for index, listener := range e.listeners {
			if listener.Id == id {
				copy(e.listeners[index:], e.listeners[index+1:])
				e.listeners = e.listeners[:len(e.listeners)-1]
				return
			}
		}
		panic(fmt.Errorf("listener not added to event"))
	}

	fn := reflect.TypeOf(signature)
	e.paramTypes = make([]reflect.Type, fn.NumIn())

	for index := range e.paramTypes {
		e.paramTypes[index] = fn.In(index)
	}

	e.isVariadic = fn.IsVariadic()
}

func (e *EventBase) String() string {
	// TODO : @Badu - use strings.Builder
	s := "Event<"
	for i, t := range e.paramTypes {
		if i > 0 {
			s += ", "
		}
		if e.isVariadic && i == len(e.paramTypes)-1 {
			s += "..."
		}
		s += t.String()
	}
	return s + ">"
}

func assignable(to, from reflect.Type) bool {
	if from == nil {
		switch to.Kind() {
		case reflect.Chan, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.Func:
			return true
		}
		return false
	}
	return from.AssignableTo(to)
}

func (e *EventBase) VerifySignature(argTys []reflect.Type, isVariadic bool) {
	if !isVariadic {
		if len(e.paramTypes) != len(argTys) {
			panic(fmt.Errorf("%v.Emit(%v) Argument count mismatch. Expected %d, got %d", e.String(), argTys, len(e.paramTypes), len(argTys)))
		}
		for i, argTy := range argTys {
			paramTy := e.paramTypes[i]
			if !assignable(paramTy, argTy) {
				panic(fmt.Errorf("%v.Emit(%v) Argument %v for was of the wrong type. Got: %v, Expected: %v", e.String(), argTys, i, argTy, paramTy))
			}
		}
		return
	}

	// variadic
	if len(argTys) < len(e.paramTypes)-1 {
		panic(fmt.Errorf("%v.Emit(%v) Too few arguments. Must have at least %v, but got %v", e.String(), argTys, len(e.paramTypes), len(argTys)))
	}

	varIdx := len(e.paramTypes) - 1
	for index, argTy := range argTys {
		if index >= varIdx {
			paramTy := e.paramTypes[varIdx].Elem()
			if !assignable(paramTy, argTy) {
				panic(fmt.Errorf("%v.Emit(%v) Variadic argument %v for was of the wrong type. Got: %v, Expected: %v", e.String(), argTys, index-varIdx, argTy, paramTy))
			}
		} else {
			paramTy := e.paramTypes[index]
			if !assignable(paramTy, argTy) {
				panic(fmt.Errorf("%v.Emit(%v) Argument %v for was of the wrong type. Got: %v, Expected: %v", e.String(), argTys, index, argTy, paramTy))
			}
		}
	}
}

func (e *EventBase) VerifyArguments(args []interface{}) {
	argTys := make([]reflect.Type, len(args))
	for i, arg := range args {
		argTys[i] = reflect.TypeOf(arg)
	}
	e.VerifySignature(argTys, e.isVariadic)
}

func (e *EventBase) InvokeListeners(args []interface{}) {
	argVals := make([]reflect.Value, len(args))
	for i, arg := range args {
		if arg == nil {
			argVals[i] = reflect.New(e.paramTypes[i]).Elem()
		} else {
			argVals[i] = reflect.ValueOf(arg)
		}
	}

	for _, listener := range e.listeners {
		listener.Function.Call(argVals)
	}
}

// Event compliance
func (e *EventBase) Listen(listener interface{}) EventSubscription {
	var paramTypes []reflect.Type
	var function reflect.Value

	reflectTy := reflect.TypeOf(listener)
	if reflectTy.Kind() == reflect.Func {
		paramTypes = make([]reflect.Type, reflectTy.NumIn())
		for i, _ := range paramTypes {
			paramTypes[i] = reflectTy.In(i)
		}
		function = reflect.ValueOf(listener)
	} else {
		switch ty := listener.(type) {
		case Event:
			paramTypes = ty.ParameterTypes()
			function = reflect.ValueOf(listener).MethodByName("Emit")
		default:
			panic(fmt.Errorf("listener cannot be of type %v", reflectTy.String()))
		}
	}

	if function.IsNil() {
		panic("listener function is nil")
	}

	for i, listenerTy := range paramTypes {
		if !listenerTy.AssignableTo(e.paramTypes[i]) {
			panic(fmt.Errorf("%v.Listen(%v) Listener parameter %v for was of the wrong type. Got: %v, Expected: %v", e.String(), listener, i, listenerTy, e.paramTypes[i]))
		}
	}

	id := e.nextId
	e.nextId++

	e.listeners = append(e.listeners, EventListener{Id: id, Function: function})

	return &eventBaseSubscription{e, id}
}

func (e *EventBase) Emit(args ...interface{}) {
	e.VerifyArguments(args)
	e.InvokeListeners(args)
}

func (e *EventBase) ParameterTypes() []reflect.Type {
	return e.paramTypes
}
