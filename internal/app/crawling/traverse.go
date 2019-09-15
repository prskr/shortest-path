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
	"net/http"
	"regexp"
)

var (
	baseDomainRegex = regexp.MustCompile(`^http(s)?://[A-z]+\.wikipedia.org`)
)

func NewWikiCrawler(startPage string, targetPath string, maxHops uint16) *WikiCrawler {
	return &WikiCrawler{
		alreadyVisitedPages: newStringSet(),
		startPage:           startPage,
		targetPage:          targetPath,
		wikiBaseDomain:      baseDomainRegex.FindString(startPage),
		maxHops:             maxHops,
	}
}

type WikiCrawler struct {
	alreadyVisitedPages *stringSet
	startPage           string
	targetPage          string
	wikiBaseDomain      string
	maxHops             uint16
	fetchedPages        uint
}

func (crawler WikiCrawler) FetchedPages() uint {
	return crawler.fetchedPages
}

func (crawler WikiCrawler) DiscoveredPages() int {
	return len(crawler.alreadyVisitedPages.items)
}

func (crawler *WikiCrawler) SearchShortestPath() (traversalResult TraversalResult, err error) {
	var depth uint16 = 0

	var currentStates = []*TraversalState{{
		PageURI:     crawler.startPage,
		Predecessor: nil,
		Ancestors:   nil,
	}}

	for {
		if depth >= crawler.maxHops {
			err = fmt.Errorf("reached max hops")
			return
		}

		retrievedStates := make([]*TraversalState, 0)

		for _, s := range currentStates {
			if traversalResult = crawler.processState(s); traversalResult.foundPath() {
				return
			}
			retrievedStates = append(retrievedStates, s.Ancestors...)
		}

		currentStates = retrievedStates
		depth += 1
	}
}

func (crawler *WikiCrawler) processState(state *TraversalState) (traversalResult TraversalResult) {
	logger := log.WithFields(log.Fields{
		"pageURI": state.PageURI,
	})

	var err error
	var resp *http.Response

	logger.Debug("Fetching wiki page")

	resp, err = http.Get(state.PageURI)

	crawler.fetchedPages += 1

	if err != nil {
		logger.Error("failed to retrieve page URI")
		return
	}

	logger.Debug("Parsing retrieved HTML page")

	var discoveredLinks []string
	discoveredLinks, err = extractLinksFromContent(resp.Body, crawler.alreadyVisitedPages, func(s string) string {
		return fmt.Sprintf("%s%s", crawler.wikiBaseDomain, s)
	})

	if err != nil {
		logger.WithError(err).Errorf("Failed to process page %s", state.PageURI)
		return
	}

	for _, link := range discoveredLinks {
		ancestor := &TraversalState{
			PageURI:     link,
			Predecessor: state,
		}

		if link == crawler.targetPage {
			traversalResult.successState = ancestor
			return
		}

		state.Ancestors = append(state.Ancestors, ancestor)
	}

	return
}
