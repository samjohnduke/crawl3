package crawler

import (
	"net/url"
	"time"
)

// Crawl represents the output from fetching a webpage and parsing/extracting its
// contents. It also contains meta data about the page, timing details and
// an error if it was encounted. As a Value object it contains no methods
type Crawl struct {
	URL      string
	ID       string
	PageHash string

	LoadedTime  time.Time
	StartTime   time.Time
	FetchTime   time.Time
	ExtractTime time.Time
	EndTime     time.Time

	Title       string
	Description string

	HarvestedURLs []string
	HarvestedData interface{}
	MicroData     interface{}
	MetaData      map[string]interface{}
	JSONData      []interface{}
	RawData       string

	Error     string
	ErrorCode string

	// A signal will be sent when the crawl has been completed
	sig chan struct{}
	url *url.URL
}

// Host will return the hostname of the url of the crawl
func (c *Crawl) Host() string {
	u, err := url.Parse(c.URL)
	if err != nil {
		return ""
	}

	return u.Host
}
