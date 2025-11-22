// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"fmt"
	"reflect"

	"github.com/badu/gxui/math"
)

// AdapterItem is a user defined type that can be used to uniquely identify a
// single item in an adapter. The type must support equality and be hashable.
type AdapterItem interface{}

type Viewer interface {
	View(styles *StyleDefs) Control
}

type DefaultAdapter struct {
	AdapterBase
	itemToIndex map[AdapterItem]int
	styleLabel  func(styles *StyleDefs, label *Label)
	items       reflect.Value
	size        math.Size
}

func CreateDefaultAdapter(width, height int) *DefaultAdapter {
	return &DefaultAdapter{size: math.Size{W: width, H: height}}
}

func (a *DefaultAdapter) SetSizeAsLargest(styles *StyleDefs) {
	size := math.Size{}
	font := styles.DefaultFont
	for index := 0; index < a.Count(); index++ {
		switch t := a.ItemAt(index).(type) {
		case Viewer:
			size = size.Max(t.View(styles).DesiredSize(math.ZeroSize, math.MaxSize))

		case fmt.Stringer:
			size = size.Max(font.Measure(&TextBlock{Runes: []rune(t.String())}))

		default:
			size = size.Max(font.Measure(&TextBlock{Runes: []rune(fmt.Sprintf("%+v", t))}))
		}
	}
	a.SetSize(size)
}

func (a *DefaultAdapter) SetStyleLabel(providerFn func(styles *StyleDefs, label *Label)) {
	a.styleLabel = providerFn
	a.DataChanged(true)
}

func (a *DefaultAdapter) Count() int {
	if !a.items.IsValid() {
		return 0
	}

	switch a.items.Kind() {
	case reflect.Slice, reflect.Array:
		return a.items.Len()

	default:
		return 1
	}
}

func (a *DefaultAdapter) ItemAt(index int) AdapterItem {
	count := a.Count()
	if index < 0 || index >= count {
		panic(fmt.Errorf("ItemAt index %d is out of bounds [%d, %d]", index, 0, count-1))
	}

	switch a.items.Kind() {
	case reflect.Slice, reflect.Array:
		return a.items.Index(index).Interface()

	default:
		return a.items.Interface()
	}

}

func (a *DefaultAdapter) ItemIndex(item AdapterItem) int {
	return a.itemToIndex[item]
}

func (a *DefaultAdapter) Size(styles *StyleDefs) math.Size {
	return a.size
}

func (a *DefaultAdapter) SetSize(size math.Size) {
	a.size = size
	a.DataChanged(true)
}

func (a *DefaultAdapter) Create(driver Driver, styles *StyleDefs, index int) Control {
	switch t := a.ItemAt(index).(type) {
	case Viewer:
		return t.View(styles)

	case fmt.Stringer:
		label := CreateLabel(driver, styles)
		label.SetMargin(math.ZeroSpacing)
		label.SetMultiline(false)
		label.SetText(t.String())
		if a.styleLabel != nil {
			a.styleLabel(styles, label)
		}
		return label

	default:
		label := CreateLabel(driver, styles)
		label.SetMargin(math.ZeroSpacing)
		label.SetMultiline(false)
		label.SetText(fmt.Sprintf("%+v", t))
		if a.styleLabel != nil {
			a.styleLabel(styles, label)
		}
		return label
	}
}

func (a *DefaultAdapter) Items() interface{} {
	return a.items.Interface()
}

func (a *DefaultAdapter) SetItems(items interface{}) {
	a.items = reflect.ValueOf(items)
	a.itemToIndex = make(map[AdapterItem]int)
	for idx := 0; idx < a.Count(); idx++ {
		a.itemToIndex[a.ItemAt(idx)] = idx
	}
	a.DataReplaced()
}
