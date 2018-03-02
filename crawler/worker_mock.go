package crawler

// workerMock is the minimum code required to do nothing but still function as expected,
// useful of testing the dispatcher
type workerMock struct {
	pool chan chan *Crawl
	jobs chan *Crawl
	quit chan chan bool
}

func newWorkerMock(pool chan chan *Crawl) Worker {
	return &workerMock{
		pool: pool,
		jobs: make(chan *Crawl),
		quit: make(chan chan bool),
	}
}

func (w *workerMock) Start() error {
	go w.work()
	return nil
}

func (w *workerMock) Stop() error {
	wait := make(chan bool)
	w.quit <- wait
	<-wait
	return nil
}

func (w *workerMock) work() {
	for {
		// register the current worker into the worker queue.
		w.pool <- w.jobs

		select {
		case <-w.jobs:
			break

		case q := <-w.quit:
			q <- true
			return
		}
	}
}
