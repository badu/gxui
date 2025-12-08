package main

import (
	"time"

	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/purego"
	"github.com/badu/gxui/pkg/font"
	"github.com/badu/gxui/pkg/math"
	"github.com/badu/gxui/samples/flags"
	"github.com/chewxy/math32"
)

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)

	font, err := driver.CreateFont(font.Default, 75)
	if err != nil {
		panic(err)
	}

	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Hi")
	window.SetBackgroundBrush(gxui.CreateBrush(gxui.Gray50))

	label := gxui.CreateLabel(driver, styles)
	label.SetFont(font)
	label.SetText("Hello world")

	window.AddChild(label)

	ticker := time.NewTicker(time.Millisecond * 30)
	go func() {
		phase := float32(0)
		for _ = range ticker.C {
			redCos := math32.Cos((phase + 0.000) * math.TwoPi)
			greenCos := math32.Cos((phase + 0.333) * math.TwoPi)
			blueCos := math32.Cos((phase + 0.666) * math.TwoPi)
			alphaCos := math32.Cos(phase * 10)
			c := gxui.Color{
				R: 0.75 + 0.25*redCos,
				G: 0.75 + 0.25*greenCos,
				B: 0.75 + 0.25*blueCos,
				A: 0.50 + 0.50*alphaCos,
			}
			phase += 0.01
			driver.Call(func() {
				label.SetColor(c)
			})
		}
	}()

	window.OnClose(ticker.Stop)
	window.OnClose(driver.Terminate)
}

func main() {
	purego.StartDriver(appMain)
}
