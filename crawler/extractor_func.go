package crawler

import "github.com/PuerkitoBio/goquery"

// FuncExtractor allows you to use a Function to extract data from a webpage.
type FuncExtractor struct {
	URL    string
	Method func(doc *goquery.Document) (interface{}, error)
}

// Register the extractor into a list of extractors
func (fn *FuncExtractor) Register(ex Extractors) {
	ex.Add(fn)
}

// Match simply returns the url that this function will be run against
func (fn *FuncExtractor) Match() (url string) {
	return fn.URL
}

// Extract is passeed a goquery.Document and we use the provided function to extract some data
func (fn *FuncExtractor) Extract(doc *goquery.Document) (harvestedData interface{}, err error) {
	return fn.Method(doc)
}
