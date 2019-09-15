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

import "golang.org/x/net/html"

func newStringSet() *stringSet {
	return &stringSet{items: make(map[string]bool)}
}

type stringSet struct {
	items map[string]bool
}

func (set *stringSet) Add(value string) (alreadyPresent bool) {
	if _, alreadyPresent = set.items[value]; !alreadyPresent {
		set.items[value] = true
	}
	return
}

func (set *stringSet) Contains(value string) (contained bool) {
	_, contained = set.items[value]
	return
}

type tokenStack struct {
	tokens []html.Token
}

func (stack *tokenStack) Push(token html.Token) {
	stack.tokens = append(stack.tokens, token)
}

func (stack tokenStack) Peek() (token html.Token, ok bool) {
	if ok = len(stack.tokens) > 0; !ok {
		return
	}

	n := len(stack.tokens) - 1
	token = stack.tokens[n]
	return
}

func (stack *tokenStack) Pop() (token html.Token, ok bool) {
	if ok = len(stack.tokens) > 0; !ok {
		return
	}

	n := len(stack.tokens) - 1
	token = stack.tokens[n]
	stack.tokens = stack.tokens[:n]
	return
}

func (stack tokenStack) Empty() bool {
	return len(stack.tokens) == 0
}
