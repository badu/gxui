// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import (
	"fmt"
)

// A 3x3 matrix:
//
//	╭          ╮
//	│ M₀ M₁ M₂ │
//	│ M₃ M₄ M₅ │
//	│ M₆ M₇ M₈ │
//	╰          ╯
type Mat3 [9]float32

func CreateMat3(r0c0, r0c1, r0c2, r1c0, r1c1, r1c2, r2c0, r2c1, r2c2 float32) Mat3 {
	return Mat3{r0c0, r0c1, r0c2, r1c0, r1c1, r1c2, r2c0, r2c1, r2c2}
}

func (m Mat3) String() string {
	s := make([]string, 9)
	l := 0
	for i, v := range m {
		s[i] = fmt.Sprintf("%.5f", v)
		l = max(l, len(s[i]))
	}
	for i := range m {
		for len(s[i]) < l {
			s[i] = " " + s[i]
		}
	}
	p := ""
	for i := 0; i < l; i++ {
		p += " "
	}
	return fmt.Sprintf(
		"\n╭ %s %s %s ╮"+
			"\n│ %s %s %s │"+
			"\n│ %s %s %s │"+
			"\n│ %s %s %s │"+
			"\n╰ %s %s %s ╯",
		p, p, p,
		s[0], s[1], s[2],
		s[3], s[4], s[5],
		s[6], s[7], s[8],
		p, p, p,
	)
}
