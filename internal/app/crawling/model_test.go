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
