// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"math"
)

var ZeroSize Size
var ZeroSpacing Spacing
var ZeroPoint Point

var MaxSize = Size{Width: 0x40000000, Height: 0x40000000} // Picked to avoid integer overflows

const Pi = float32(math.Pi)
const TwoPi = Pi * 2.0
