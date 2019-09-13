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
