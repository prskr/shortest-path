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
	"reflect"
	"testing"
)

func Test_stringSet_Add(t *testing.T) {
	type fields struct {
		items map[string]bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantAlreadyPresent bool
		resultingSetSize   int
	}{
		{
			name: "Test add into empty set",
			fields: fields{
				items: make(map[string]bool),
			},
			args: args{
				value: "hello, world",
			},
			wantAlreadyPresent: false,
			resultingSetSize:   1,
		},
		{
			name: "Test add already existing string",
			fields: fields{
				items: map[string]bool{
					"hello, world": true,
				},
			},
			args: args{
				value: "hello, world",
			},
			wantAlreadyPresent: true,
			resultingSetSize:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := &stringSet{
				items: tt.fields.items,
			}
			if gotAlreadyPresent := set.Add(tt.args.value); gotAlreadyPresent != tt.wantAlreadyPresent {
				t.Errorf("Add() = %v, want %v", gotAlreadyPresent, tt.wantAlreadyPresent)
			}

			if tt.resultingSetSize != len(set.items) {
				t.Errorf("Expected resulting set to have size %d but got %d", tt.resultingSetSize, len(set.items))
			}
		})
	}
}

func Test_stringSet_Contains(t *testing.T) {
	type fields struct {
		items map[string]bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantContained bool
	}{
		{
			name: "check for value in empty set",
			fields: fields{
				items: make(map[string]bool),
			},
			args: args{
				value: "hello, world",
			},
			wantContained: false,
		},
		{
			name: "check for actually contained value",
			fields: fields{
				items: map[string]bool{
					"hello, world": true,
				},
			},
			args: args{
				value: "hello, world",
			},
			wantContained: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := &stringSet{
				items: tt.fields.items,
			}
			if gotContained := set.Contains(tt.args.value); gotContained != tt.wantContained {
				t.Errorf("Contains() = %v, want %v", gotContained, tt.wantContained)
			}
		})
	}
}

func Test_tokenStack_Push(t *testing.T) {
	type fields struct {
		tokens []html.Token
	}
	type args struct {
		token html.Token
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expectedLength int
	}{
		{
			name: "Push to empty stack",
			fields: fields{
				make([]html.Token, 0),
			},
			args: args{
				token: html.Token{
					Type:     html.StartTagToken,
					DataAtom: 0,
					Data:     "",
					Attr: []html.Attribute{
						html.Attribute{
							Namespace: "a",
							Key:       "href",
							Val:       "/wiki/Times_New_Roman",
						},
					},
				},
			},
			expectedLength: 1,
		},
		{
			name: "Push to not empty stack",
			fields: fields{
				[]html.Token{
					html.Token{
						Type:     html.StartTagToken,
						DataAtom: 0,
						Data:     "a",
						Attr: []html.Attribute{
							html.Attribute{
								Namespace: "",
								Key:       "href",
								Val:       "/wiki/Times_New_Roman",
							},
						},
					},
				},
			},
			args: args{
				token: html.Token{
					Type:     html.StartTagToken,
					DataAtom: 0,
					Data:     "b",
					Attr:     make([]html.Attribute, 0),
				},
			},
			expectedLength: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &tokenStack{
				tokens: tt.fields.tokens,
			}
			stack.Push(tt.args.token)
			if len(stack.tokens) != tt.expectedLength {
				t.Errorf("expected length %d but got %d", tt.expectedLength, len(stack.tokens))
			}
		})
	}
}

func Test_tokenStack_Pop(t *testing.T) {
	type fields struct {
		tokens []html.Token
	}
	tests := []struct {
		name      string
		fields    fields
		wantToken html.Token
		wantOk    bool
	}{
		{
			name: "Pop from not empty stack",
			fields: fields{
				tokens: []html.Token{
					html.Token{
						Type:     html.StartTagToken,
						DataAtom: 0,
						Data:     "b",
						Attr:     make([]html.Attribute, 0),
					},
				},
			},
			wantToken: html.Token{
				Type:     html.StartTagToken,
				DataAtom: 0,
				Data:     "b",
				Attr:     make([]html.Attribute, 0),
			},
			wantOk: true,
		},
		{
			name: "Pop from empty stack",
			fields: fields{
				tokens: make([]html.Token, 0),
			},
			wantToken: html.Token{},
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &tokenStack{
				tokens: tt.fields.tokens,
			}
			gotToken, gotOk := stack.Pop()
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("Pop() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Pop() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_tokenStack_Peek(t *testing.T) {
	type fields struct {
		tokens []html.Token
	}
	tests := []struct {
		name           string
		fields         fields
		wantToken      html.Token
		wantOk         bool
		expectedLength int
	}{
		{
			name: "Peek from not empty stack",
			fields: fields{
				tokens: []html.Token{
					html.Token{
						Type:     html.StartTagToken,
						DataAtom: 0,
						Data:     "b",
						Attr:     make([]html.Attribute, 0),
					},
				},
			},
			wantToken: html.Token{
				Type:     html.StartTagToken,
				DataAtom: 0,
				Data:     "b",
				Attr:     make([]html.Attribute, 0),
			},
			wantOk:         true,
			expectedLength: 1,
		},
		{
			name: "Peek from empty stack",
			fields: fields{
				tokens: make([]html.Token, 0),
			},
			wantToken:      html.Token{},
			wantOk:         false,
			expectedLength: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tokenStack{
				tokens: tt.fields.tokens,
			}
			gotToken, gotOk := stack.Peek()
			if !reflect.DeepEqual(gotToken, tt.wantToken) {
				t.Errorf("Peek() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Peek() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if len(stack.tokens) != tt.expectedLength {
				t.Errorf("Expected length %d but got %d", tt.expectedLength, len(stack.tokens))
			}
		})
	}
}
