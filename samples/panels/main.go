// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/gl"
	"github.com/badu/gxui/samples/flags"
)

// Create a PanelHolderImpl with a 3 panels
func panelHolder(name string, driver gxui.Driver, styles *gxui.StyleDefs) gxui.PanelHolder {
	label := func(text string) gxui.Label {
		label := gxui.CreateLabel(driver, styles)
		label.SetText(text)
		return label
	}

	holder := gxui.CreatePanelHolder(driver, styles)
	holder.AddPanel(label(name+" 0 content"), name+" 0 panel")
	holder.AddPanel(label(name+" 1 content"), name+" 1 panel")
	holder.AddPanel(label(name+" 2 content"), name+" 2 panel")
	return holder
}

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	// ┌───────┐║┌───────┐
	// │       │║│       │
	// │   A   │║│   B   │
	// │       │║│       │
	// └───────┘║└───────┘
	// ═══════════════════
	// ┌───────┐║┌───────┐
	// │       │║│       │
	// │   C   │║│   D   │
	// │       │║│       │
	// └───────┘║└───────┘

	splitterAB := gxui.CreateSplitterLayout(driver, styles)
	splitterAB.SetOrientation(gxui.Horizontal)
	splitterAB.AddChild(panelHolder("A", driver, styles))
	splitterAB.AddChild(panelHolder("B", driver, styles))

	splitterCD := gxui.CreateSplitterLayout(driver, styles)
	splitterCD.SetOrientation(gxui.Horizontal)
	splitterCD.AddChild(panelHolder("C", driver, styles))
	splitterCD.AddChild(panelHolder("D", driver, styles))

	vSplitter := gxui.CreateSplitterLayout(driver, styles)
	vSplitter.SetOrientation(gxui.Vertical)
	vSplitter.AddChild(splitterAB)
	vSplitter.AddChild(splitterCD)

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Panels")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(vSplitter)
	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
