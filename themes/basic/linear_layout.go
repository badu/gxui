// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basic

import (
	"github.com/badu/gxui"
)

func CreateLinearLayout(theme *Theme) gxui.LinearLayout {
	l := &gxui.LinearLayoutImpl{}
	l.Init(l, theme)
	return l
}
