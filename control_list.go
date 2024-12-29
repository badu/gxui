// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type ControlList []Control

func (l ControlList) Contains(control Control) bool {
	for _, item := range l {
		if item == control {
			return true
		}
	}

	return false
}
