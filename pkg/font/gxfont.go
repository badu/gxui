package font

//go:generate go run mkfont.go

import (
	"bytes"
	"compress/flate"
	"io"
)

var (
	// Default is the standard GXUI sans-serif font.
	Default = inflate(robotoRegular)

	// Monospace is the standard GXUI fixed-width font.
	Monospace = inflate(droidSansMono)
)

func inflate(src []byte) []byte {
	reader := bytes.NewReader(src)
	data, err := io.ReadAll(flate.NewReader(reader))
	if err != nil {
		panic(err)
	}
	return data
}
