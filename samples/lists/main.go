// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/gl"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/samples/flags"
)

// Number picker uses the gxui.DefaultAdapter for driving a list
func numberPicker(driver gxui.Driver, styles *gxui.StyleDefs, overlay gxui.BubbleOverlay) gxui.Control {
	items := []string{
		"zero", "one", "two", "three", "four", "five",
		"six", "seven", "eight", "nine", "ten",
		"eleven", "twelve", "thirteen", "fourteen", "fifteen",
		"sixteen", "seventeen", "eighteen", "nineteen", "twenty",
	}

	adapter := gxui.CreateDefaultAdapter()
	adapter.SetItems(items)

	layout := gxui.CreateLinearLayout(driver, styles)
	layout.SetDirection(gxui.TopToBottom)

	label0 := gxui.CreateLabel(driver, styles)
	label0.SetText("Numbers:")
	layout.AddChild(label0)

	dropList := gxui.CreateDropDownList(driver, styles)
	dropList.SetAdapter(adapter)
	dropList.SetBubbleOverlay(overlay)
	layout.AddChild(dropList)

	list := gxui.CreateList(driver, styles)
	list.SetAdapter(adapter)
	list.SetOrientation(gxui.Vertical)
	layout.AddChild(list)

	label1 := gxui.CreateLabel(driver, styles)
	label1.SetMargin(math.Spacing{T: 30})
	label1.SetText("Selected number:")
	layout.AddChild(label1)

	selected := gxui.CreateLabel(driver, styles)
	layout.AddChild(selected)

	dropList.OnSelectionChanged(func(item gxui.AdapterItem) {
		if list.Selected() != item {
			list.Select(item)
		}
	})

	list.OnSelectionChanged(func(item gxui.AdapterItem) {
		if dropList.Selected() != item {
			dropList.Select(item)
		}
		selected.SetText(fmt.Sprintf("%s - %d", item, adapter.ItemIndex(item)))
	})

	return layout
}

type customAdapter struct {
	gxui.AdapterBase
}

func (a *customAdapter) Count() int {
	return 1000
}

func (a *customAdapter) ItemAt(index int) gxui.AdapterItem {
	return index // This adapter uses integer indices as AdapterItems
}

func (a *customAdapter) ItemIndex(item gxui.AdapterItem) int {
	return item.(int) // Inverse of ItemAt()
}

func (a *customAdapter) Size(styles *gxui.StyleDefs) math.Size {
	return math.Size{W: 100, H: 100}
}

func (a *customAdapter) Create(driver gxui.Driver, styles *gxui.StyleDefs, index int) gxui.Control {
	phase := float32(index) / 1000
	c := gxui.Color{
		R: 0.5 + 0.5*math.Sinf(math.TwoPi*(phase+0.000)),
		G: 0.5 + 0.5*math.Sinf(math.TwoPi*(phase+0.333)),
		B: 0.5 + 0.5*math.Sinf(math.TwoPi*(phase+0.666)),
		A: 1.0,
	}
	i := gxui.CreateImage(driver, styles)
	i.SetBackgroundBrush(gxui.CreateBrush(c))
	i.SetMargin(math.Spacing{L: 3, T: 3, R: 3, B: 3})
	i.OnMouseEnter(func(ev gxui.MouseEvent) {
		i.SetBorderPen(gxui.CreatePen(2, gxui.Gray80))
	})
	i.OnMouseExit(func(ev gxui.MouseEvent) {
		i.SetBorderPen(gxui.TransparentPen)
	})
	i.OnMouseDown(func(ev gxui.MouseEvent) {
		i.SetBackgroundBrush(gxui.CreateBrush(c.MulRGB(0.7)))
	})
	i.OnMouseUp(func(ev gxui.MouseEvent) {
		i.SetBackgroundBrush(gxui.CreateBrush(c))
	})
	return i
}

// Color picker uses the customAdapter for driving a list
func colorPicker(driver gxui.Driver, styles *gxui.StyleDefs) gxui.Control {
	layout := gxui.CreateLinearLayout(driver, styles)
	layout.SetDirection(gxui.TopToBottom)

	label0 := gxui.CreateLabel(driver, styles)
	label0.SetText("Color palette:")
	layout.AddChild(label0)

	adapter := &customAdapter{}

	list := gxui.CreateList(driver, styles)
	list.SetAdapter(adapter)
	list.SetOrientation(gxui.Horizontal)
	layout.AddChild(list)

	label1 := gxui.CreateLabel(driver, styles)
	label1.SetMargin(math.Spacing{T: 30})
	label1.SetText("Selected color:")
	layout.AddChild(label1)

	selected := gxui.CreateImage(driver, styles)
	selected.SetExplicitSize(math.Size{W: 32, H: styles.FontSize + 8})
	layout.AddChild(selected)

	list.OnSelectionChanged(func(item gxui.AdapterItem) {
		if item != nil {
			control := list.ItemControl(item)
			selected.SetBackgroundBrush(control.(gxui.Image).BackgroundBrush())
		}
	})

	return layout
}

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	overlay := gxui.CreateBubbleOverlay(driver, styles)

	holder := gxui.CreatePanelHolder(driver, styles)
	holder.AddPanel(numberPicker(driver, styles, overlay), "Default adapter")
	holder.AddPanel(colorPicker(driver, styles), "Custom adapter")

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Lists")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(holder)
	window.AddChild(overlay)
	window.OnClose(driver.Terminate)
	window.SetPadding(math.Spacing{L: 10, T: 10, R: 10, B: 10})
}

func main() {
	gl.StartDriver(appMain)
}
