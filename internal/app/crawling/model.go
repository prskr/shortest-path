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

type TraversalState struct {
	PageURI     string
	Predecessor *TraversalState
	Ancestors   []*TraversalState
}

type TraversalResult struct {
	successState *TraversalState
}

func (tr TraversalResult) foundPath() bool {
	return tr.successState != nil
}

func (tr TraversalResult) VisitedPages() (visitedPages []string) {
	currentState := tr.successState

	for currentState != nil {
		visitedPages = append(visitedPages, currentState.PageURI)
		currentState = currentState.Predecessor
	}
	return
}
