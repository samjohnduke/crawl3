package crawler

import "context"

type dispatcher struct {
	workerCount int64
	workerQueue chan chan *Crawl
	workQueue   chan *Crawl
	workers     []Worker
	quit        chan chan bool
	newWorker   WorkerFactoryFunc
}

func newDispatcher(count int64, queue chan *Crawl, createWorkerFunc WorkerFactoryFunc) *dispatcher {
	dispatcher := &dispatcher{
		workerCount: count,
		workerQueue: make(chan chan *Crawl, count),
		workQueue:   queue,
		workers:     []Worker{},
		quit:        make(chan chan bool),
		newWorker:   createWorkerFunc,
	}

	return dispatcher
}

func (d *dispatcher) Stop(ctx *context.Context) error {
	wait := make(chan bool)
	d.quit <- wait

	<-wait

	return nil
}

func (d *dispatcher) Start() error {
	var i int64
	for i = 0; i < d.workerCount; i++ {
		worker := d.newWorker(d.workerQueue)
		d.workers = append(d.workers, worker)
		worker.Start()
	}

	go d.Dispatcher()
	return nil
}

func (d *dispatcher) Dispatcher() {
	for {
		select {
		case job := <-d.workQueue:
			// a job request has been received
			jobChannel := <-d.workerQueue
			jobChannel <- job
			break

		case q := <-d.quit:
			for i := range d.workers {
				d.workers[i].Stop()
			}

			q <- true
			return
		}
	}
}
