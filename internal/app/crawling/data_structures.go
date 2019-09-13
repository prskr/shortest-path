package crawling

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