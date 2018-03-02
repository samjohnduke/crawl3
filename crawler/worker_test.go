package crawler

import (
	"log"
	"os"
	"testing"
)

func TestWorker(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	ins := NewInstrumentationMem()
	exes := newDefaultExtractors()

	opts := WorkerOpts{
		logger:     logger,
		instrument: ins,
		extractors: exes,
		results:    make(chan *Crawl),
	}

	worker := NewDefaultWorker(nil, opts)
	for _, u := range testWorkerUrls {
		crawl := &Crawl{
			URL: u,
		}
		err := worker.(*defaultWorker).do(crawl)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInvalidURLWorker(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	ins := NewInstrumentationMem()
	exes := newDefaultExtractors()

	opts := WorkerOpts{
		logger:     logger,
		instrument: ins,
		extractors: exes,
		results:    make(chan *Crawl),
	}

	worker := NewDefaultWorker(nil, opts)

	for _, u := range testWorkerInvalidURL {
		crawl := &Crawl{
			URL: u,
		}
		err := worker.(*defaultWorker).do(crawl)
		if err == nil {
			t.Error(err)
		}
	}
}

var testWorkerUrls = []string{
	"https://www.domain.com.au/bourke-street-melbourne-vic-3000-10828929",
}

var testWorkerInvalidURL = []string{
	"test.local",
	"http://test.local",
	"ftp://test.local",
}
