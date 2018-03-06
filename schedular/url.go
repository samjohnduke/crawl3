package schedular

import (
	"net/url"
	"sync"

	"github.com/PuerkitoBio/purell"
)

// A URL is the high level representation of a POSSIBLE url. It takes a url string and
// performs some basic operations on it to ensure that its in its most basic form
type URL struct {
	original   string
	normalised string
	url        *url.URL
}

// NewURL parses the url as a string and normalises it for later deduplication
func NewURL(u string) (*URL, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	retURL := &URL{
		original: u,
		url:      ur,
	}

	retURL.normalised, err = purell.NormalizeURLString(u, purell.FlagsSafe|purell.FlagRemoveFragment)
	if err != nil {
		return nil, err
	}

	return retURL, nil
}

// NewURLWithReference returns a new URL using the reference as a guide
func NewURLWithReference(u string, ref string) (*URL, error) {
	retURL, err := NewURL(u)

	refURL, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}

	err = retURL.ResolveReference(refURL)
	if err != nil {
		return nil, err
	}

	return retURL, nil
}

// ResolveReference uses the provided url as a guide to determin what the
// url should have that its missing (ie relative to absolute)
func (u *URL) ResolveReference(root *url.URL) error {
	u.url = root.ResolveReference(u.url)

	var err error
	u.normalised, err = purell.NormalizeURLString(u.url.String(), purell.FlagsSafe|purell.FlagRemoveFragment)
	if err != nil {
		return err
	}

	return nil
}

type urlList struct {
	list []*URL
	mu   sync.Mutex
}

func newURLList() *urlList {
	return &urlList{
		list: []*URL{},
	}
}

func (ul *urlList) Len() int {
	return len(ul.list)
}

func (ul *urlList) push(newURL *URL) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.list = append(ul.list, newURL)
}

func (ul *urlList) pop(c int) []*URL {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	p := ul.list[len(ul.list)-c-1:]
	ul.list = ul.list[:len(ul.list)-c-1]
	return p
}

func (ul *urlList) shift(c int) []*URL {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	p := ul.list[0:c]
	ul.list = ul.list[c:]

	return p
}

func (ul *urlList) unshift(newURL *URL) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.list = append([]*URL{newURL}, ul.list...)
}
