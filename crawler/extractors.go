package crawler

// The Extractors is an object that builds a list of possible extractors.
// the implementation will perform optimization based what the extractor's
// listen for (thoughts for the future)
type extractors interface {
	Add(extractor)
	Matches(url string) []extractor
}

// defaultExtractors creates a mapping from a url to an extractor
type defaultExtractors struct {
	list map[string][]extractor
}

func newDefaultExtractors() extractors {
	return &defaultExtractors{
		list: make(map[string][]extractor),
	}
}

func (e *defaultExtractors) Add(ex extractor) {
	url := ex.Match()
	list := e.list[url]

	if list == nil {
		l := make([]extractor, 0)
		l = append(l, ex)
		e.list[url] = l
	} else {
		e.list[url] = append(e.list[url], ex)
	}
}

func (e *defaultExtractors) Matches(url string) []extractor {
	return e.list[url]
}
