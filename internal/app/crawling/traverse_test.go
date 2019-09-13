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
	"testing"
)

func Test_processState(t *testing.T) {
	type args struct {
		crawler *WikiCrawler
		state   *TraversalState
	}
	tests := []struct {
		name             string
		args             args
		wantResultsCount int
	}{
		{
			name: "Test to fetch Times 'New Roman article'",
			args: args{
				crawler: NewWikiCrawler("https://en.wikipedia.org/wiki/Times_New_Roman", "https://en.wikipedia.org/wiki/Great_Britain", 10),
				state: &TraversalState{
					PageURI:     "https://en.wikipedia.org/wiki/Times_New_Roman",
					Predecessor: nil,
				},
			},
			wantResultsCount: 352,
		},
		{
			name: "Test to fetch image link",
			args: args{
				crawler: NewWikiCrawler("https://en.wikipedia.org/wiki/Times_New_Roman", "https://en.wikipedia.org/wiki/Great_Britain", 10),
				state: &TraversalState{
					PageURI:     "https://en.wikipedia.org/wiki/File:Times_New_Roman-sample.svg",
					Predecessor: nil,
				},
			},
			wantResultsCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.args.crawler.processState(tt.args.state)

			if len(tt.args.state.Ancestors) != tt.wantResultsCount {
				t.Errorf("expected %d results but got %d", tt.wantResultsCount, len(tt.args.state.Ancestors))
			}
		})
	}
}
