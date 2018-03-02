package crawler

import "time"

// Crawl represents the output from fetching a webpage and parsing/extracting its
// contents. It also contains meta data about the page, timing details and
// an error if it was encounted. As a Value object it contains no methods
type Crawl struct {
	URL string
	ID  string

	LoadedTime  time.Time
	StartTime   time.Time
	FetchTime   time.Time
	ExtractTime time.Time
	EndTime     time.Time

	HarvestedURLs []string
	HarvestedData interface{}
	MicroData     interface{}
	MetaData      map[string]string
	JSONData      []interface{}
	RawData       string

	Error error

	// A signal will be sent when the crawl has been completed
	sig chan struct{}
}
