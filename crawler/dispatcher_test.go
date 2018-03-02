package crawler

import (
	"testing"
)

// This test will detect contention of queuing channels
func TestSingleWorkerDispatcher(t *testing.T) {
	queue := make(chan *Crawl)

	workerFactory := func(pool chan chan *Crawl) Worker {
		return newWorkerMock(pool)
	}

	dispatcher := newDispatcher(1, queue, workerFactory)

	err := dispatcher.Start()
	if err != nil {
		t.Error(err)
	}

	queue <- &Crawl{URL: "not-a-link"}
	queue <- &Crawl{URL: "also-not-a-link"}
	queue <- &Crawl{URL: "another-not-link"}

	err = dispatcher.Stop(nil)
	if err != nil {
		t.Error(err)
	}
}

func TestMultiWorkerDispatcher(t *testing.T) {
	queue := make(chan *Crawl)

	workerFactory := func(pool chan chan *Crawl) Worker {
		return newWorkerMock(pool)
	}

	dispatcher := newDispatcher(4, queue, workerFactory)

	err := dispatcher.Start()
	if err != nil {
		t.Error(err)
	}

	queue <- &Crawl{URL: "not-a-link"}
	queue <- &Crawl{URL: "also-not-a-link"}
	queue <- &Crawl{URL: "another-not-link"}
	queue <- &Crawl{URL: "another-not-link"}
	queue <- &Crawl{URL: "another-not-link"}
	queue <- &Crawl{URL: "another-not-link"}
	queue <- &Crawl{URL: "another-not-link"}
	queue <- &Crawl{URL: "another-not-link"}

	err = dispatcher.Stop(nil)
	if err != nil {
		t.Error(err)
	}
}
