package crawler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/iand/microdata"
)

// Worker is the abstract interface for crawling a webpage
type Worker interface {
	Start() error
	Stop() error
}

type defaultWorker struct {
	pool       chan chan *Crawl
	jobs       chan *Crawl
	quit       chan chan bool
	results    chan *Crawl
	instrument Instrument
	logger     *log.Logger
	extractors extractors
}

type WorkerFactoryFunc func(chan chan *Crawl) Worker
type WorkerFactoryInvoker func(WorkerOpts) WorkerFactoryFunc

// WorkerOpts are the outside components that are required by the worker
// to do its job or observe its behaviour
type WorkerOpts struct {
	logger     *log.Logger
	instrument Instrument
	extractors extractors
	results    chan *Crawl
}

// NewDefaultWorker creates a new worker based on the options provided
func NewDefaultWorker(pool chan chan *Crawl, opts WorkerOpts) Worker {
	return &defaultWorker{
		pool:       pool,
		results:    opts.results,
		jobs:       make(chan *Crawl),
		quit:       make(chan chan bool),
		instrument: opts.instrument,
		logger:     opts.logger,
		extractors: opts.extractors,
	}
}

func (w *defaultWorker) Start() error {
	go w.work()
	return nil
}

func (w *defaultWorker) Stop() error {
	wait := make(chan bool)
	w.quit <- wait

	<-wait

	return nil
}

func (w *defaultWorker) work() {
	for {
		// register the current worker into the worker queue.
		w.pool <- w.jobs

		select {
		case job := <-w.jobs:
			w.do(job)
			//w.results <- job

			var m struct{}
			job.sig <- m

			break

		case q := <-w.quit:
			q <- true
			return
		}
	}
}

func (w *defaultWorker) do(u *Crawl) error {
	w.instrument.Gauge("workers_active", 1)
	u.StartTime = time.Now()

	parsedURL, err := url.Parse(u.URL)
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err
		return err
	}

	resp, err := http.Get(u.URL)
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err
		return err
	}

	u.FetchTime = time.Now()

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err
		return err
	}

	urls := doc.Find("a[href]").Map(func(i int, sel *goquery.Selection) string {
		if attr, ok := sel.Attr("href"); ok {
			return attr
		}
		return ""
	})

	metadata := make(map[string]string)
	doc.Find("meta[name]").Each(func(_ int, sel *goquery.Selection) {
		var name string
		var value string

		if attr, ok := sel.Attr("name"); ok {
			name = attr
		}

		if attr, ok := sel.Attr("content"); ok {
			value = attr
		}

		if name != "" && value != "" {
			metadata[name] = value
		}
	})
	doc.Find("meta[property]").Each(func(_ int, sel *goquery.Selection) {
		var name string
		var value string

		if attr, ok := sel.Attr("property"); ok {
			name = attr
		}

		if attr, ok := sel.Attr("content"); ok {
			value = attr
		}

		if name != "" && value != "" {
			metadata[name] = value
		}
	})

	jld := []interface{}{}
	doc.Find("script[type='application/ld+json']").Each(func(_ int, sel *goquery.Selection) {
		content := sel.Text()
		var parsedContent interface{}
		err := json.Unmarshal([]byte(content), &parsedContent)
		if err != nil {
			return
		}
		jld = append(jld, parsedContent)
	})

	ps := microdata.NewParser(bytes.NewReader(body), parsedURL)
	mdata, err := ps.Parse()
	if err != nil {
		w.logger.Println(err)
		u.Error = err
	}

	var harvested []interface{}
	fncs := w.extractors.Matches(parsedURL.Host)
	for _, fn := range fncs {
		h, err := fn.Extract(doc)
		if err != nil {
			w.logger.Println(err)
			continue
		}

		harvested = append(harvested, h)
	}

	u.ExtractTime = time.Now()

	u.MicroData = mdata
	u.HarvestedData = harvested
	u.HarvestedURLs = urls
	u.JSONData = jld
	u.MetaData = metadata

	w.instrument.Gauge("workers_active", -1)
	w.instrument.Count("crawl_url")

	u.EndTime = time.Now()
	return nil
}
