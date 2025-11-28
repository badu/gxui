package main

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/cgo"
	"github.com/badu/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	label1 := gxui.CreateLabel(driver, styles)
	label1.SetColor(gxui.White)
	label1.SetText("1x1")

	cell1x1 := gxui.CreateLinearLayout(driver, styles)
	cell1x1.SetBackgroundBrush(gxui.CreateBrush(gxui.Blue40))
	cell1x1.SetHorizontalAlignment(gxui.AlignCenter)
	cell1x1.AddChild(label1)

	label2 := gxui.CreateLabel(driver, styles)
	label2.SetColor(gxui.White)
	label2.SetText("2x1")

	cell2x1 := gxui.CreateLinearLayout(driver, styles)
	cell2x1.SetBackgroundBrush(gxui.CreateBrush(gxui.Green40))
	cell2x1.SetHorizontalAlignment(gxui.AlignCenter)
	cell2x1.AddChild(label2)

	label3 := gxui.CreateLabel(driver, styles)
	label3.SetColor(gxui.White)
	label3.SetText("1x2")

	cell1x2 := gxui.CreateLinearLayout(driver, styles)
	cell1x2.SetBackgroundBrush(gxui.CreateBrush(gxui.Red40))
	cell1x2.SetHorizontalAlignment(gxui.AlignCenter)
	cell1x2.AddChild(label3)

	table := gxui.CreateTableLayout(driver, styles)
	table.SetGrid(3, 2) // columns, rows

	// row, column, horizontal span, vertical span
	table.SetChildAt(0, 0, 1, 1, cell1x1)
	table.SetChildAt(0, 1, 2, 1, cell2x1)
	table.SetChildAt(2, 0, 1, 2, cell1x2)

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Table")
	window.AddChild(table)
	window.OnClose(driver.Terminate)
}

func main() {
	cgo.StartDriver(appMain)
}
