package shared

import "sync"

// A URLList is a concurrently safe list of URLS
type URLList struct {
	list []*URL
	mu   *sync.Mutex
}

// NewURLList creates a new URLList using a slice and a sync mutex
func NewURLList() *URLList {
	return &URLList{
		list: []*URL{},
		mu:   &sync.Mutex{},
	}
}

// Len gets the number of URLS in the list
func (ul *URLList) Len() int {
	return len(ul.list)
}

// Push a new url onto the end of the list
func (ul *URLList) Push(newURL *URL) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.list = append(ul.list, newURL)
}

// Pop c (count) urls from the beginning of the list
func (ul *URLList) Pop(c int) []*URL {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	p := ul.list[len(ul.list)-c:]
	ul.list = ul.list[:len(ul.list)-c]
	return p
}

// Shift pops c (count) urls from the end of the list
func (ul *URLList) Shift(c int) []*URL {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	p := ul.list[0:c]
	ul.list = ul.list[c:]

	return p
}

// Unshift pushes a url onto the beginning of the list
func (ul *URLList) Unshift(newURL *URL) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.list = append([]*URL{newURL}, ul.list...)
}
