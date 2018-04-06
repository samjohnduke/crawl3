package crawler

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

// The Transport is the interface for recieving and sending crawls over the
// implmeneted methods
type Transport interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// CrawlRequest Sends a request to the service to crawl a page
type CrawlRequest struct {
	URL string
}

// CrawlAsyncRequest sends a request with a reply option
type CrawlAsyncRequest struct {
	URL   string
	Reply string
}

// ProgressRequest asks the service how the crawl is progressing
type ProgressRequest struct {
	GUID string
}

// CrawlReply - All requests return a crawl and an option error if something went wrong
type CrawlReply struct {
	Crawl Crawl
	Error error
}

// The Service is the interface that must be implemented to communicate
// with this service. It provides a means for receiving
type Service interface {
	// Crawl a url and return the result
	Crawl(ctx context.Context, url string) (result *Crawl, err error)

	// Crawl a result and return a guid that points to an in progress crawl
	CrawlAsync(ctx context.Context, url string, cb func(*Crawl)) (guid string, err error)

	// Get the progress of a crawl.
	CrawlProgress(ctx context.Context, guid string) (result *Crawl, err error)
}

// The Publisher is an interface that will push out the result of a crawl to those
// who want ictx context.Contextt
type Publisher interface {
	Publish(crawl *Crawl) error
}

// A Listener provides a way for the spider to get updates on its crawls in the system
// It is the other half of Publisher
type Listener interface {
	Listen() chan *Crawl
	Close() error
}

// The Client interface is the
type Client interface {
	//An asynchronous request for crawling a webpage
	CrawlAsync(ctx context.Context, url string, cb func(*Crawl)) (guid string, err error)

	//A synchronous request for crawling a webpage
	Crawl(ctx context.Context, url string) (result *Crawl, err error)

	// Get the progress of a crawl
	CrawlProgress(ctx context.Context, guid string) (result *Crawl, err error)
}

// An Instrument is a way for the application to monitor the state of the running application
// and how it has behaved in the past. By watching the data we can take action. This interface
// allows multiple different Instrumentations to be used (in memory or prometheus for example)
type Instrument interface {
	Count(metric string)
	Gauge(metric string, value int64)
	Histogram(metric string, value string)
}

// Crawler responds to the service for performing a crawl
type crawler struct {
	Timeout     time.Duration
	Client      *http.Client
	Concurrency int64
	Queue       chan *Crawl
	Open        map[string]*Crawl
	Dispatcher  *dispatcher
	Output      chan *Crawl
	Lock        sync.RWMutex
}

// ServiceOpts are optional interfaces that are used through the system
type ServiceOpts struct {
	Logger      *log.Logger
	Instrument  Instrument
	Extractors  Extractors
	Publisher   Publisher
	WorkerCount int64
}

// New creates the core service that will be used to crawl with
func New(opts ServiceOpts, workerFactoryInv WorkerFactoryInvoker) (Service, error) {
	output := make(chan *Crawl)
	queue := make(chan *Crawl)

	var ins Instrument
	var logger *log.Logger
	var exes Extractors
	var factory WorkerFactoryFunc

	if opts.Instrument == nil {
		ins = NewInstrumentationMem()
	} else {
		ins = opts.Instrument
	}

	if opts.Logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		logger = opts.Logger
	}

	if opts.Extractors == nil {
		exes = NewDefaultExtractors()
	} else {
		exes = opts.Extractors
	}

	workerOpts := WorkerOpts{
		logger:     logger,
		extractors: exes,
		instrument: ins,
		results:    output,
		publisher:  opts.Publisher,
	}

	if workerFactoryInv == nil {
		// The worker factory uses a closure to enable adding higher level components to a lower
		// level piece of functionality and hiding the implementation of the dispatcher
		factory = func(pool chan chan *Crawl) Worker {
			return NewDefaultWorker(pool, workerOpts)
		}
	} else {
		factory = workerFactoryInv(workerOpts)
	}

	dispatcher := newDispatcher(opts.WorkerCount, queue, factory)
	err := dispatcher.Start()
	if err != nil {
		return nil, err
	}

	return &crawler{
		Timeout:     10 * time.Second,
		Client:      http.DefaultClient,
		Concurrency: 4,
		Queue:       queue,
		Open:        make(map[string]*Crawl),
		Dispatcher:  dispatcher,
		Output:      output,
	}, nil
}

//An asynchronous request for crawling a webpage
func (c *crawler) CrawlAsync(ctx context.Context, url string, cb func(*Crawl)) (guid string, err error) {
	crawl := &Crawl{
		URL: url,
		sig: make(chan struct{}),
	}
	gd, err := uuid.NewV4()
	crawl.ID = gd.String()

	crawl.LoadedTime = time.Now()

	c.loadCrawl(crawl)

	c.Queue <- crawl

	go func() {
		<-crawl.sig
		c.unloadCrawl(crawl)
		cb(crawl)
	}()

	return crawl.ID, err
}

//A synchronous request for crawling a webpage
func (c *crawler) Crawl(ctx context.Context, url string) (result *Crawl, err error) {
	log.Println(url)
	crawl := &Crawl{
		URL: url,
		sig: make(chan struct{}),
	}
	gd, err := uuid.NewV4()
	crawl.ID = gd.String()

	crawl.LoadedTime = time.Now()

	c.loadCrawl(crawl)

	c.Queue <- crawl

	<-crawl.sig
	c.unloadCrawl(crawl)

	return crawl, err
}

// Get the progress of a crawl
func (c *crawler) CrawlProgress(ctx context.Context, guid string) (result *Crawl, err error) {
	return nil, nil
}

func (c *crawler) getCrawl(guid string) *Crawl {
	c.Lock.RLock()
	crawl := c.Open[guid]
	c.Lock.RUnlock()
	return crawl
}

func (c *crawler) loadCrawl(crawl *Crawl) {
	c.Lock.Lock()
	c.Open[crawl.ID] = crawl
	c.Lock.Unlock()
}

func (c *crawler) unloadCrawl(crawl *Crawl) {
	c.Lock.Lock()
	delete(c.Open, crawl.ID)
	c.Lock.Unlock()
}
