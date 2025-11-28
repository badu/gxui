package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/badu/gxui"
	"github.com/badu/gxui/drivers/cgo"
	"github.com/badu/gxui/samples/flags"
)

func appMain(driver gxui.Driver) {
	args := flag.Args()
	file := ""
	if len(args) != 1 {
		fmt.Print("usage: image_viewer image-path. Not provided, using default\n")
		file = "./samples/image_viewer/sasha_mishu.jpg"
	} else {
		file = args[0]
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("Failed to open image '%s': %v\n", file, err)
		os.Exit(1)
	}

	source, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("Failed to read image '%s': %v\n", file, err)
		os.Exit(1)
	}

	styles := flags.CreateTheme(driver)
	img := gxui.CreateImage(driver, styles)

	mx := source.Bounds().Max
	window := gxui.CreateWindow(driver, styles, mx.X, mx.Y, "Image viewer")
	window.SetScale(flags.DefaultScaleFactor)
	window.AddChild(img)

	// Copy the image to a RGBA format before handing to a gxui.Texture
	rgba := image.NewRGBA(source.Bounds())
	draw.Draw(rgba, source.Bounds(), source, image.ZP, draw.Src)
	texture := driver.CreateTexture(rgba, 1)
	img.SetTexture(texture)

	window.OnClose(driver.Terminate)
}

func main() {
	flag.Parse()
	cgo.StartDriver(appMain)
}
