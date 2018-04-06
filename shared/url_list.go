package shared

import "sync"

type URLList struct {
	list []*URL
	mu   *sync.Mutex
}

func NewURLList() *URLList {
	return &URLList{
		list: []*URL{},
		mu:   &sync.Mutex{},
	}
}

func (ul *URLList) Len() int {
	return len(ul.list)
}

func (ul *URLList) push(newURL *URL) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.list = append(ul.list, newURL)
}

func (ul *URLList) pop(c int) []*URL {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	p := ul.list[len(ul.list)-c:]
	ul.list = ul.list[:len(ul.list)-c]
	return p
}

func (ul *URLList) shift(c int) []*URL {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	p := ul.list[0:c]
	ul.list = ul.list[c:]

	return p
}

func (ul *URLList) unshift(newURL *URL) {
	ul.mu.Lock()
	defer ul.mu.Unlock()

	ul.list = append([]*URL{newURL}, ul.list...)
}
