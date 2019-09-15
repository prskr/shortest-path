// Copyright © 2019 Peter Kurfer peter.kurfer@googlemail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crawling

import (
	"reflect"
	"testing"
)

func TestTraversalResult_VisitedPages(t *testing.T) {
	type fields struct {
		latestState *TraversalState
	}
	tests := []struct {
		name             string
		fields           fields
		wantVisitedPages []string
	}{
		{
			name: "get visited pages if previous state is nil",
			fields: fields{
				latestState: nil,
			},
			wantVisitedPages: nil,
		},
		{
			name: "get visited pages for single previous state",
			fields: fields{
				latestState: &TraversalState{
					PageURI:     "https://en.wikipedia.org/wiki/Times_New_Roman",
					Predecessor: nil,
					Ancestors:   nil,
				},
			},
			wantVisitedPages: []string{"https://en.wikipedia.org/wiki/Times_New_Roman"},
		},
		{
			name: "get visited pages for a graph of pages",
			fields: fields{
				latestState: &TraversalState{
					PageURI: "https://en.wikipedia.org/wiki/The_Times",
					Predecessor: &TraversalState{
						PageURI:     "https://en.wikipedia.org/wiki/Times_New_Roman",
						Predecessor: nil,
						Ancestors:   nil,
					},
					Ancestors: nil,
				},
			},
			wantVisitedPages: []string{
				"https://en.wikipedia.org/wiki/The_Times",
				"https://en.wikipedia.org/wiki/Times_New_Roman",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TraversalResult{
				successState: tt.fields.latestState,
			}
			if gotVisitedPages := tr.VisitedPages(); !reflect.DeepEqual(gotVisitedPages, tt.wantVisitedPages) {
				t.Errorf("VisitedPages() = %v, want %v", gotVisitedPages, tt.wantVisitedPages)
			}
		})
	}
}
