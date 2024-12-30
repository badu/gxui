// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package flags holds command line options common to all GXUI samples.
package flags

import (
	"flag"
	"fmt"
	"github.com/badu/gxui"
	"github.com/badu/gxui/gxfont"
	"github.com/goxjs/glfw"
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
func CreateTheme(driver gxui.Driver) *gxui.StyleDefs {
	if FlagTheme == "light" {
		return CreateLightTheme(driver, FontSize)
	}
	return CreateDarkTheme(driver, FontSize)
}

func CreateLightTheme(driver gxui.Driver, fontSize int) *gxui.StyleDefs {
	defaultFont, err := driver.CreateFont(gxfont.Default, fontSize)
	if err == nil {
		defaultFont.LoadGlyphs(32, 126)
	} else {
		fmt.Printf("Warning: Failed to load default font - %v\n", err)
	}

	defaultMonospaceFont, err := driver.CreateFont(gxfont.Monospace, fontSize)
	if err == nil {
		defaultFont.LoadGlyphs(32, 126)
	} else {
		fmt.Printf("Warning: Failed to load default monospace font - %v\n", err)
	}

	scrollBarRailDefaultBg := gxui.Black
	scrollBarRailDefaultBg.A = 0.7

	scrollBarRailOverBg := gxui.Gray20
	scrollBarRailOverBg.A = 0.7

	neonBlue := gxui.ColorFromHex(0xFF5C8CFF)
	focus := gxui.ColorFromHex(0xFFC4D6FF)

	monitor := glfw.GetPrimaryMonitor()
	_, _, w, h := monitor.GetWorkarea()
	if w == 0 || h == 0 {
		vm := monitor.GetVideoMode()
		w, h = vm.Width, vm.Height
	}

	return &gxui.StyleDefs{
		DefaultFont:          defaultFont,
		DefaultMonospaceFont: defaultMonospaceFont,
		WindowBackground:     gxui.White,

		//                                   fontColor    brushColor   penColor
		BubbleOverlayStyle:        gxui.CreateStyle(gxui.Gray40, gxui.Gray20, gxui.Gray40, 1.0),
		ButtonDefaultStyle:        gxui.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray40, 1.0),
		ButtonOverStyle:           gxui.CreateStyle(gxui.Gray40, gxui.Gray90, gxui.Gray40, 1.0),
		ButtonPressedStyle:        gxui.CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		CodeSuggestionListStyle:   gxui.CreateStyle(gxui.Gray40, gxui.Gray20, gxui.Gray10, 1.0),
		DropDownListDefaultStyle:  gxui.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray20, 1.0),
		DropDownListOverStyle:     gxui.CreateStyle(gxui.Gray40, gxui.Gray90, gxui.Gray50, 1.0),
		FocusedStyle:              gxui.CreateStyle(gxui.Gray20, gxui.Transparent, focus, 1.0),
		HighlightStyle:            gxui.CreateStyle(gxui.Gray40, gxui.Transparent, neonBlue, 2.0),
		LabelStyle:                gxui.CreateStyle(gxui.Gray40, gxui.Transparent, gxui.Transparent, 0.0),
		PanelBackgroundStyle:      gxui.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray15, 1.0),
		ScrollBarBarDefaultStyle:  gxui.CreateStyle(gxui.Gray40, gxui.Gray30, gxui.Gray40, 1.0),
		ScrollBarBarOverStyle:     gxui.CreateStyle(gxui.Gray40, gxui.Gray50, gxui.Gray60, 1.0),
		ScrollBarRailDefaultStyle: gxui.CreateStyle(gxui.Gray40, scrollBarRailDefaultBg, gxui.Transparent, 1.0),
		ScrollBarRailOverStyle:    gxui.CreateStyle(gxui.Gray40, scrollBarRailOverBg, gxui.Gray20, 1.0),
		SplitterBarDefaultStyle:   gxui.CreateStyle(gxui.Gray40, gxui.Gray80, gxui.Gray40, 1.0),
		SplitterBarOverStyle:      gxui.CreateStyle(gxui.Gray40, gxui.Gray80, gxui.Gray50, 1.0),
		TabActiveHighlightStyle:   gxui.CreateStyle(gxui.Gray30, neonBlue, neonBlue, 0.0),
		TabDefaultStyle:           gxui.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray40, 1.0),
		TabOverStyle:              gxui.CreateStyle(gxui.Gray30, gxui.Gray90, gxui.Gray50, 1.0),
		TabPressedStyle:           gxui.CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		TextBoxDefaultStyle:       gxui.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray20, 1.0),
		TextBoxOverStyle:          gxui.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray50, 1.0),

		ScreenWidth:  w,
		ScreenHeight: h,
		FontSize:     fontSize,
	}
}

func CreateDarkTheme(driver gxui.Driver, fontSize int) *gxui.StyleDefs {
	defaultFont, err := driver.CreateFont(gxfont.Default, fontSize)
	if err == nil {
		defaultFont.LoadGlyphs(32, 126)
	} else {
		fmt.Printf("Warning: Failed to load default font - %v\n", err)
	}

	defaultMonospaceFont, err := driver.CreateFont(gxfont.Monospace, fontSize)
	if err == nil {
		defaultFont.LoadGlyphs(32, 126)
	} else {
		fmt.Printf("Warning: Failed to load default monospace font - %v\n", err)
	}

	scrollBarRailDefaultBg := gxui.Black
	scrollBarRailDefaultBg.A = 0.7

	scrollBarRailOverBg := gxui.Gray20
	scrollBarRailOverBg.A = 0.7

	neonBlue := gxui.ColorFromHex(0xFF5C8CFF)
	focus := gxui.ColorFromHex(0xA0C4D6FF)

	monitor := glfw.GetPrimaryMonitor()
	_, _, w, h := monitor.GetWorkarea()
	if w == 0 || h == 0 {
		vm := monitor.GetVideoMode()
		w, h = vm.Width, vm.Height
	}

	return &gxui.StyleDefs{
		DefaultFont:          defaultFont,
		DefaultMonospaceFont: defaultMonospaceFont,
		WindowBackground:     gxui.Black,

		//                                   fontColor    brushColor   penColor
		BubbleOverlayStyle:        gxui.CreateStyle(gxui.Gray80, gxui.Gray20, gxui.Gray40, 1.0),
		ButtonDefaultStyle:        gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray20, 1.0),
		ButtonOverStyle:           gxui.CreateStyle(gxui.Gray90, gxui.Gray15, gxui.Gray50, 1.0),
		ButtonPressedStyle:        gxui.CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		CodeSuggestionListStyle:   gxui.CreateStyle(gxui.Gray80, gxui.Gray20, gxui.Gray10, 1.0),
		DropDownListDefaultStyle:  gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray20, 1.0),
		DropDownListOverStyle:     gxui.CreateStyle(gxui.Gray80, gxui.Gray15, gxui.Gray50, 1.0),
		FocusedStyle:              gxui.CreateStyle(gxui.Gray80, gxui.Transparent, focus, 1.0),
		HighlightStyle:            gxui.CreateStyle(gxui.Gray80, gxui.Transparent, neonBlue, 2.0),
		LabelStyle:                gxui.CreateStyle(gxui.Gray80, gxui.Transparent, gxui.Transparent, 0.0),
		PanelBackgroundStyle:      gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray15, 1.0),
		ScrollBarBarDefaultStyle:  gxui.CreateStyle(gxui.Gray80, gxui.Gray30, gxui.Gray40, 1.0),
		ScrollBarBarOverStyle:     gxui.CreateStyle(gxui.Gray80, gxui.Gray50, gxui.Gray60, 1.0),
		ScrollBarRailDefaultStyle: gxui.CreateStyle(gxui.Gray80, scrollBarRailDefaultBg, gxui.Transparent, 1.0),
		ScrollBarRailOverStyle:    gxui.CreateStyle(gxui.Gray80, scrollBarRailOverBg, gxui.Gray20, 1.0),
		SplitterBarDefaultStyle:   gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray10, 1.0),
		SplitterBarOverStyle:      gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray50, 1.0),
		TabActiveHighlightStyle:   gxui.CreateStyle(gxui.Gray90, neonBlue, neonBlue, 0.0),
		TabDefaultStyle:           gxui.CreateStyle(gxui.Gray80, gxui.Gray30, gxui.Gray40, 1.0),
		TabOverStyle:              gxui.CreateStyle(gxui.Gray90, gxui.Gray30, gxui.Gray50, 1.0),
		TabPressedStyle:           gxui.CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		TextBoxDefaultStyle:       gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray20, 1.0),
		TextBoxOverStyle:          gxui.CreateStyle(gxui.Gray80, gxui.Gray10, gxui.Gray50, 1.0),

		ScreenWidth:  w,
		ScreenHeight: h,
		FontSize:     fontSize,
	}
}
