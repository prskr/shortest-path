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
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"regexp"
)

var (
	wikiLinkRegex = regexp.MustCompile(`^/wiki/[A-z_\-#()]+$`)
)

type linkFormatter func(string) string

func extractLinksFromContent(body io.ReadCloser, visitedPages *stringSet, formatter linkFormatter) (links []string, err error) {
	tokenStack := tokenStack{}
	tokenizer := html.NewTokenizer(body)

	var token html.Token
	token, err = seekDOMElementBySelector(tokenizer, "div", "id", "bodyContent")

	if err != nil {
		return
	}

	tokenStack.Push(token)

	for !tokenStack.Empty() && tokenizer.Next() != html.ErrorToken {
		currentToken := tokenizer.Token()
		switch currentToken.Type {
		case html.EndTagToken:
			latestToken, ok := tokenStack.Peek()
			if !ok {
				err = fmt.Errorf("token stack already empty")
				return
			}
			if latestToken.Data != currentToken.Data {
				err = fmt.Errorf("latest token on stack does not match closing tag")
				return
			} else {
				tokenStack.Pop()
			}
			break
		case html.StartTagToken, html.SelfClosingTagToken:

			// push tag to stack to be able to track closing tags
			if currentToken.Type == html.StartTagToken {
				tokenStack.Push(currentToken)
			}
			if currentToken.Data == "a" {
				for _, attr := range currentToken.Attr {
					if requireAll(
						attr.Key == "href",
						wikiLinkRegex.MatchString(attr.Val),
						!visitedPages.Contains(attr.Val),
					) {
						log.Debugf("Enqueuing discovered link %s", attr.Val)
						visitedPages.Add(attr.Val)
						links = append(links, formatter(attr.Val))
					}
				}
			}
			break
		}
	}
	return
}

func seekDOMElementBySelector(tokenizer *html.Tokenizer, elementType, selectorKey, selectorValue string) (token html.Token, err error) {
	for tokenizer.Next() != html.ErrorToken {
		switch token = tokenizer.Token(); token.Type {
		case html.StartTagToken, html.SelfClosingTagToken:
			if token.Data == elementType {
				for _, attr := range token.Attr {
					if attr.Key == selectorKey && attr.Val == selectorValue {
						return
					}
				}
			}
		}
	}
	err = fmt.Errorf("requested element with type %s and selector %s=%s not found", elementType, selectorKey, selectorValue)
	return
}

func requireAll(booleans ...bool) bool {
	result := true

	for _, b := range booleans {
		result = result && b
		if !result {
			return result
		}
	}

	return result
}
