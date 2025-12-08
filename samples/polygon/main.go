package main

import (
	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/purego"
	"github.com/badu/gxui/pkg/math"
	"github.com/badu/gxui/samples/flags"
	"github.com/chewxy/math32"
)

func drawStar(canvas gxui.Canvas, center math.Point, radius, rotation float32, points int) {
	p := make(gxui.Polygon, points*2)
	for i := 0; i < points*2; i++ {
		frac := float32(i) / float32(points*2)
		α := frac*math.TwoPi + rotation
		r := []float32{radius, radius / 2}[i&1]
		sinα, cosα := math32.Sincos(α)
		p[i] = gxui.PolygonVertex{
			Position: math.Point{
				X: center.X + int(r*cosα),
				Y: center.Y + int(r*sinα),
			},
			RoundedRadius: []float32{0, 50}[i&1],
		}
	}
	canvas.DrawPolygon(p, gxui.CreatePen(3, gxui.Red), gxui.CreateBrush(gxui.Yellow))
}

func drawMoon(canvas gxui.Canvas, center math.Point, radius float32) {
	c := 40
	p := make(gxui.Polygon, c*2)
	for i := 0; i < c; i++ {
		frac := float32(i) / float32(c)
		α := math.Lerpf(math.Pi*1.2, math.Pi*-0.2, frac)
		sinα, cosα := math32.Sincos(α)
		p[i] = gxui.PolygonVertex{
			Position: math.Point{
				X: center.X + int(radius*sinα),
				Y: center.Y + int(radius*cosα),
			},
			RoundedRadius: 0,
		}
	}

	for i := 0; i < c; i++ {
		frac := float32(i) / float32(c)
		α := math.Lerpf(math.Pi*-0.2, math.Pi*1.2, frac)
		r := math.Lerpf(radius, radius*0.5, math32.Sin(frac*math.Pi))
		sinα, cosα := math32.Sincos(α)
		p[i+c] = gxui.PolygonVertex{
			Position: math.Point{
				X: center.X + int(r*sinα),
				Y: center.Y + int(r*cosα),
			},
			RoundedRadius: 0,
		}
	}
	canvas.DrawPolygon(p, gxui.CreatePen(3, gxui.Gray80), gxui.CreateBrush(gxui.Gray40))
}

func appMain(driver gxui.Driver) {
	styles := flags.CreateTheme(driver)
	window := gxui.CreateWindow(driver, styles, styles.ScreenWidth/2, styles.ScreenHeight, "Polygon")
	window.SetScale(flags.DefaultScaleFactor)

	canvas := driver.CreateCanvas(math.Size{Width: 1000, Height: 1000})
	drawStar(canvas, math.Point{X: 100, Y: 100}, 50, 0.2, 6)
	drawStar(canvas, math.Point{X: 650, Y: 170}, 70, 0.5, 7)
	drawStar(canvas, math.Point{X: 40, Y: 300}, 20, 0, 5)
	drawStar(canvas, math.Point{X: 410, Y: 320}, 25, 0.9, 5)
	drawStar(canvas, math.Point{X: 220, Y: 520}, 45, 0, 6)

	drawMoon(canvas, math.Point{X: 400, Y: 300}, 200)
	canvas.Complete()

	image := gxui.CreateImage(driver, styles)
	image.SetCanvas(canvas)
	window.AddChild(image)

	window.OnClose(driver.Terminate)
}

func main() {
	purego.StartDriver(appMain)
}
