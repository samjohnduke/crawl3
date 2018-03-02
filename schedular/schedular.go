package schedular

import (
	"net/url"
)

// Service is the interface for launching the spider against a particular service border
// and the implmentation should use this appropriately
type Service interface {
	Schedule(rootURL string) error
}

// The Schedular is responsible for managing the application state of a crawl. It should
// listen for crawls, pull out the harvested urls, check for whether or not it is needed
// to crawl and then if yes intelligently push them to the crawler
type Schedular struct {
}

// A URL is the high level representation of a POSSIBLE url. It takes a url string and
// performs some basic operations on it to ensure that its in its most basic form
type URL struct {
	original   string
	normalised string
	url        *url.URL
}
