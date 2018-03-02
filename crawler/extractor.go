package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
)

func init() {

}

// Extractor pulls data from a page
type extractor interface {
	Register(extractors)
	Match() (url string)
	Extract(doc *goquery.Document) (harvestedData interface{}, err error)
}

type javaScriptExtractor struct {
	js string
	vm otto.Otto
}

func (jse *javaScriptExtractor) Register(ex extractors) {
	ex.Add(jse)
}

func (jse *javaScriptExtractor) Match() (url string) {
	return ""
}

func (jse *javaScriptExtractor) Extract(doc *goquery.Document) (harvestedData interface{}, err error) {
	return nil, nil
}

type funcExtractor struct {
	url    string
	method func(doc *goquery.Document) (interface{}, error)
}

func (fn *funcExtractor) Register(ex extractors) {
	ex.Add(fn)
}

func (fn *funcExtractor) Match() (url string) {
	return fn.url
}

func (fn *funcExtractor) Extract(doc *goquery.Document) (harvestedData interface{}, err error) {
	return fn.method(doc)
}
