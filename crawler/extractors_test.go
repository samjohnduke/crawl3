package crawler

import (
	"net/url"
	"testing"
)

func TestExtractors(t *testing.T) {
	execs := NewDefaultExtractors()
	for _, e := range funcExs {
		e.Register(execs)
	}

	u := "http://www.abc.net.au/news/2018-02-23/barnaby-joyce-resigns/9477942"
	ur, _ := url.Parse(u)
	fncs := execs.Matches(ur.Host)
	if len(fncs) < 1 {
		t.Fail()
	}
}
