// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/gl"
	"github.com/badu/gxui/pkg/math"
	"github.com/badu/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	label := gxui.CreateLabel(driver, styles)
	label.SetText("This is a progress bar:")

	progressBar := gxui.CreateProgressBar(driver, styles)
	progressBar.SetDesiredSize(math.Size{Width: 400, Height: styles.FontSize + 4})
	progressBar.SetTarget(100)

	layout := gxui.CreateLinearLayout(driver, styles)
	layout.AddChild(label)
	layout.AddChild(progressBar)
	layout.SetHorizontalAlignment(gxui.AlignCenter)

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Progress bar")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(layout)
	window.OnClose(driver.Terminate)

	progress := 0
	pause := time.Millisecond * 500
	var timer *time.Timer
	timer = time.AfterFunc(pause, func() {
		driver.Call(func() {
			progress = (progress + 3) % progressBar.Target()
			progressBar.SetProgress(progress)
			timer.Reset(pause)
		})
	})
}

func main() {
	gl.StartDriver(appMain)
}
