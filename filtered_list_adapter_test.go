// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

import (
	"github.com/badu/gxui/test_helper"
	"strings"
	"testing"
)

func calcScore(str, partial string) int {
	return flaScore(str, strings.ToLower(str), partial, strings.ToLower(partial))
}

func TestFilteredListAdapterScore(t *testing.T) {
	test_helper.AssertEquals(t, (3*3)+(2*2+1)+(3*3+1), calcScore("a_Mix_Of_Words", "mixOfrds"))
}
