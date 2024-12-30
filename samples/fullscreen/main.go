// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/gl"
	"github.com/badu/gxui/math"
	"github.com/badu/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "WindowImpl")
	window.OnClose(driver.Terminate)
	window.SetScale(flags.DefaultScaleFactor)
	window.SetPadding(math.Spacing{L: 10, R: 10, T: 10, B: 10})
	button := gxui.CreateButton(driver, styles)
	button.SetHorizontalAlignment(gxui.AlignCenter)
	button.SetSizeMode(gxui.Fill)
	toggle := func() {
		fullscreen := !window.Fullscreen()
		window.SetFullscreen(fullscreen)
		if fullscreen {
			button.SetText("Make windowed")
		} else {
			button.SetText("Make fullscreen")
		}
	}
	button.SetText("Make fullscreen")
	button.OnClick(func(gxui.MouseEvent) { toggle() })
	window.AddChild(button)
}

func main() {
	gl.StartDriver(appMain)
}
