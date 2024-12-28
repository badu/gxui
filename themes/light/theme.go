// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package light

import (
	"fmt"
	"github.com/goxjs/glfw"

	"github.com/badu/gxui"
	"github.com/badu/gxui/gxfont"
	"github.com/badu/gxui/themes/basic"
)

func CreateTheme(driver gxui.Driver, fontSize int) gxui.Theme {
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

	return &basic.Theme{
		DriverInfo:               driver,
		DefaultFontInfo:          defaultFont,
		DefaultMonospaceFontInfo: defaultMonospaceFont,
		WindowBackground:         gxui.White,

		//                                   fontColor    brushColor   penColor
		BubbleOverlayStyle:        basic.CreateStyle(gxui.Gray40, gxui.Gray20, gxui.Gray40, 1.0),
		ButtonDefaultStyle:        basic.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray40, 1.0),
		ButtonOverStyle:           basic.CreateStyle(gxui.Gray40, gxui.Gray90, gxui.Gray40, 1.0),
		ButtonPressedStyle:        basic.CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		CodeSuggestionListStyle:   basic.CreateStyle(gxui.Gray40, gxui.Gray20, gxui.Gray10, 1.0),
		DropDownListDefaultStyle:  basic.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray20, 1.0),
		DropDownListOverStyle:     basic.CreateStyle(gxui.Gray40, gxui.Gray90, gxui.Gray50, 1.0),
		FocusedStyle:              basic.CreateStyle(gxui.Gray20, gxui.Transparent, focus, 1.0),
		HighlightStyle:            basic.CreateStyle(gxui.Gray40, gxui.Transparent, neonBlue, 2.0),
		LabelStyle:                basic.CreateStyle(gxui.Gray40, gxui.Transparent, gxui.Transparent, 0.0),
		PanelBackgroundStyle:      basic.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray15, 1.0),
		ScrollBarBarDefaultStyle:  basic.CreateStyle(gxui.Gray40, gxui.Gray30, gxui.Gray40, 1.0),
		ScrollBarBarOverStyle:     basic.CreateStyle(gxui.Gray40, gxui.Gray50, gxui.Gray60, 1.0),
		ScrollBarRailDefaultStyle: basic.CreateStyle(gxui.Gray40, scrollBarRailDefaultBg, gxui.Transparent, 1.0),
		ScrollBarRailOverStyle:    basic.CreateStyle(gxui.Gray40, scrollBarRailOverBg, gxui.Gray20, 1.0),
		SplitterBarDefaultStyle:   basic.CreateStyle(gxui.Gray40, gxui.Gray80, gxui.Gray40, 1.0),
		SplitterBarOverStyle:      basic.CreateStyle(gxui.Gray40, gxui.Gray80, gxui.Gray50, 1.0),
		TabActiveHighlightStyle:   basic.CreateStyle(gxui.Gray30, neonBlue, neonBlue, 0.0),
		TabDefaultStyle:           basic.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray40, 1.0),
		TabOverStyle:              basic.CreateStyle(gxui.Gray30, gxui.Gray90, gxui.Gray50, 1.0),
		TabPressedStyle:           basic.CreateStyle(gxui.Gray20, gxui.Gray70, gxui.Gray30, 1.0),
		TextBoxDefaultStyle:       basic.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray20, 1.0),
		TextBoxOverStyle:          basic.CreateStyle(gxui.Gray40, gxui.White, gxui.Gray50, 1.0),

		ScreenWidth:  w,
		ScreenHeight: h,
		FontSize:     fontSize,
	}
}
