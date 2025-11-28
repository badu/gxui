package purego

import (
	"image"

	"github.com/badu/gxui/pkg/math"
)

type textureContext struct {
	contextResource
	sizePixels math.Size
	texture    uint32
	flipY      bool
	pma        bool
}

type TextureImpl struct {
	image        image.Image
	pixelsPerDip float32
	flipY        bool
}

func NewTexture(fromImage image.Image, pixelsPerDip float32) *TextureImpl {
	result := &TextureImpl{
		image:        fromImage,
		pixelsPerDip: pixelsPerDip,
	}
	return result
}

// Image is gxui.Texture compliance
func (t *TextureImpl) Image() image.Image {
	return t.image
}

func (t *TextureImpl) Size() math.Size {
	return t.SizePixels().ScaleS(1.0 / t.pixelsPerDip)
}

func (t *TextureImpl) SizePixels() math.Size {
	s := t.image.Bounds().Size()
	return math.Size{Width: s.X, Height: s.Y}
}

func (t *TextureImpl) FlipY() bool {
	return t.flipY
}

func (t *TextureImpl) SetFlipY(flipY bool) {
	t.flipY = flipY
}

func (t *TextureImpl) newContext(fn *Functions) *textureContext {
	var format Enum
	var data []byte
	var pma bool

	switch imageType := t.image.(type) {
	case *image.RGBA:
		format = RGBA
		data = imageType.Pix
		pma = true
	case *image.NRGBA:
		format = RGBA
		data = imageType.Pix
	case *image.Alpha:
		format = ALPHA
		data = imageType.Pix
	default:
		panic("Unsupported image type")
	}

	glTexture := fn.CreateTexture()
	fn.BindTexture(TEXTURE_2D, glTexture)
	w, h := t.SizePixels().WH()
	fn.TexImage2D(TEXTURE_2D, 0, int32(w), int32(h), format, UNSIGNED_BYTE, data)
	fn.TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)
	fn.TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR)
	fn.BindTexture(TEXTURE_2D, uint32(0))
	checkError(fn)

	globalStats.textureContextCount.inc()
	return &textureContext{
		texture:    glTexture,
		sizePixels: t.Size(),
		flipY:      t.flipY,
		pma:        pma,
	}
}

func (c *textureContext) destroy(fn *Functions) {
	globalStats.textureContextCount.dec()
	fn.DeleteTexture(c.texture)
	c.texture = 0
}
