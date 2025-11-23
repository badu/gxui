// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

func Round(v float32) int {
	if v < 0 {
		return int(v - 0.5)
	} else {
		return int(v + 0.4999999)
	}
}

func Lerp(a, b int, s float32) int {
	r := float32(b - a)
	return a + int(r*s)
}

func Lerpf(a, b float32, s float32) float32 {
	r := b - a
	return a + r*s
}

func RampSat(s float32, a, b float32) float32 {
	return Saturate((s - a) / (b - a))
}

func Saturate(x float32) float32 {
	return Clampf(x, 0, 1)
}

func Clamp(x, min, max int) int {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

func Clampf(x, min, max float32) float32 {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}

func Mod(a, b int) int {
	x := a % b
	if x < 0 {
		return x + b
	} else {
		return x
	}
}
