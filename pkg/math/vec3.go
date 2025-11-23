// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

type Vec3 struct {
	X float32
	Y float32
	Z float32
}

func (v Vec3) Dot(vec3 Vec3) float32 {
	return v.X*vec3.X + v.Y*vec3.Y + v.Z*vec3.Z
}
