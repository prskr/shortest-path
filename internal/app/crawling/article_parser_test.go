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
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"os"
	"reflect"
	"testing"
)

func Test_extractLinksFromContent(t *testing.T) {
	type args struct {
		body         io.ReadCloser
		visitedPages *stringSet
		formatter    linkFormatter
	}
	tests := []struct {
		name            string
		args            args
		wantNumberLinks int
		wantErr         bool
	}{
		{
			name: "Get links from Manduca Jordani article",
			args: args{
				body:         MustOpen("../../../assets/test-data/manduca_jordani_article.html"),
				visitedPages: newStringSet(),
				formatter: func(s string) string {
					return s
				},
			},
			wantNumberLinks: 18,
			wantErr:         false,
		},
		{
			name: "Get links from Times New Roman article",
			args: args{
				body:         MustOpen("../../../assets/test-data/times_new_roman_article.html"),
				visitedPages: newStringSet(),
				formatter: func(s string) string {
					return s
				},
			},
			wantNumberLinks: 334,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLinks, err := extractLinksFromContent(tt.args.body, tt.args.visitedPages, tt.args.formatter)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractLinksFromContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotLinks) != tt.wantNumberLinks {
				t.Errorf("extractLinksFromContent() gotLinks = %v, want number of links %d", gotLinks, tt.wantNumberLinks)
			}
		})
	}
}

func Test_seekDOMElementBySelector(t *testing.T) {
	type args struct {
		elementType   string
		selectorKey   string
		selectorValue string
	}
	tests := []struct {
		name               string
		args               args
		wantErr            bool
		expectedChildToken html.Token
	}{
		{
			name: "seek bodyContent div",
			args: args{
				elementType:   "div",
				selectorKey:   "id",
				selectorValue: "bodyContent",
			},
			wantErr: false,
			expectedChildToken: html.Token{
				Type:     html.StartTagToken,
				Data:     "div",
				DataAtom: atom.Lookup([]byte("div")),
				Attr: []html.Attribute{
					{
						Key: "id",
						Val: "siteSub",
					},
					{
						Key: "class",
						Val: "noprint",
					},
				},
			},
		},
		{
			name: "seek content div",
			args: args{
				elementType:   "div",
				selectorKey:   "id",
				selectorValue: "content",
			},
			wantErr: false,
			expectedChildToken: html.Token{
				Type:     html.StartTagToken,
				Data:     "a",
				DataAtom: atom.Lookup([]byte("a")),
				Attr: []html.Attribute{
					{
						Key: "id",
						Val: "top",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		f, _ := os.Open("../../../assets/test-data/times_new_roman_article.html")

		tokenizer := html.NewTokenizer(f)

		t.Run(tt.name, func(t *testing.T) {
			token, err := seekDOMElementBySelector(tokenizer, tt.args.elementType, tt.args.selectorKey, tt.args.selectorValue)

			if (err != nil) != tt.wantErr {
				t.Errorf("Did not want error but got one %v", err)
			}

			if token.Data != tt.args.elementType {
				t.Errorf("Expected token to be a %s element but is %s", tt.args.elementType, token.Data)
			}

			for tokenizer.Next() == html.TextToken {
			}
			if childToken := tokenizer.Token(); !reflect.DeepEqual(childToken, tt.expectedChildToken) {
				t.Errorf("Expected child token to be %v but got %v", tt.expectedChildToken, childToken)
			}
		})
	}
}

func MustOpen(fileName string) (file *os.File) {
	file, _ = os.Open(fileName)
	return
}
