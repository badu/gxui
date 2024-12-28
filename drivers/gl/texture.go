// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl

import (
	"image"

	"github.com/badu/gxui/math"
	"github.com/goxjs/gl"
)

type texture struct {
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func newTexture(fromImage image.Image, pixelsPerDip float32) *texture {
	result := &texture{
		image:        fromImage,
		pixelsPerDip: pixelsPerDip,
	}
	return result
}

// gxui.Texture compliance
func (t *texture) Image() image.Image {
	return t.image
}

func (t *texture) Size() math.Size {
	return t.SizePixels().ScaleS(1.0 / t.pixelsPerDip)
}

func (t *texture) SizePixels() math.Size {
	s := t.image.Bounds().Size()
	return math.Size{W: s.X, H: s.Y}
}

func (t *texture) FlipY() bool {
	return t.flipY
}

func (t *texture) SetFlipY(flipY bool) {
	t.flipY = flipY
}

func (t *texture) newContext() *textureContext {
	var format gl.Enum
	var data []byte
	var pma bool

	switch imageType := t.image.(type) {
	case *image.RGBA:
		format = gl.RGBA
		data = imageType.Pix
		pma = true
	case *image.NRGBA:
		format = gl.RGBA
		data = imageType.Pix
	case *image.Alpha:
		format = gl.ALPHA
		data = imageType.Pix
	default:
		panic("Unsupported image type")
	}

	glTexture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, glTexture)
	w, h := t.SizePixels().WH()
	gl.TexImage2D(gl.TEXTURE_2D, 0, w, h, format, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture{})
	checkError()

	globalStats.textureContextCount.inc()
	return &textureContext{
		texture:    glTexture,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

type textureContext struct {
	contextResource
	texture    gl.Texture
	sizePixels math.Size
	flipY      bool
	pma        bool
}

func (c *textureContext) destroy() {
	globalStats.textureContextCount.dec()
	gl.DeleteTexture(c.texture)
	c.texture = gl.Texture{}
}
