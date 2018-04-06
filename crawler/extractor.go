package crawler

import (
	"github.com/PuerkitoBio/goquery"
)

func init() {

}

// extractor pulls data from a page
type extractor interface {
	Register(Extractors)
	Match() (url string)
	Extract(doc *goquery.Document) (harvestedData interface{}, err error)
}
