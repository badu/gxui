// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gxui

type SuggestionAdapter struct {
	FilteredListAdapter
}

func (a *SuggestionAdapter) SetSuggestions(suggestions []CodeSuggestion) {
	items := make([]FilteredListItem, len(suggestions))
	for i, s := range suggestions {
		items[i].Name = s.Name()
		items[i].Data = s
	}
	a.SetItems(items)
}

func (a *SuggestionAdapter) Suggestion(item AdapterItem) CodeSuggestion {
	return item.(FilteredListItem).Data.(CodeSuggestion)
}
