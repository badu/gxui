// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flags holds command line options common to all GXUI samples.
package flags

import (
	"flag"
	"github.com/badu/gxui"
	"github.com/badu/gxui/themes/dark"
	"github.com/badu/gxui/themes/light"
	"strconv"
)

var DefaultScaleFactor float32
var FlagTheme string
var FontSize int

func init() {
	flagTheme := flag.String("theme", "dark", "Theme to use {dark|light}.")
	fontSize := flag.String("fontSize", "24", "Adjust the font size")
	defaultScaleFactor := flag.Float64("scaling", 1.0, "Adjusts the scaling of UI rendering")
	flag.Parse()

	DefaultScaleFactor = float32(*defaultScaleFactor)
	FlagTheme = *flagTheme
	FontSize, _ = strconv.Atoi(*fontSize)
}

// CreateTheme creates and returns the theme specified on the command line.
// The default theme is dark.
func CreateTheme(driver gxui.Driver) gxui.Theme {
	if FlagTheme == "light" {
		return light.CreateTheme(driver, FontSize)
	}
	return dark.CreateTheme(driver, FontSize)
}
