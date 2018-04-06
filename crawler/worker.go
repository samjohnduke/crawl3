package crawler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/purell"
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
	extractors Extractors
	publisher  Publisher
}

// WorkerFactoryFunc is a function that takes a chan chan crawl and returns a worker
// that can be instantitated to generate workers on demand
type WorkerFactoryFunc func(chan chan *Crawl) Worker

// WorkerFactoryInvoker will return the worker factory func configured to use the
// provided opts
type WorkerFactoryInvoker func(WorkerOpts) WorkerFactoryFunc

// WorkerOpts are the outside components that are required by the worker
// to do its job or observe its behaviour
type WorkerOpts struct {
	logger     *log.Logger
	instrument Instrument
	extractors Extractors
	publisher  Publisher
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
		publisher:  opts.publisher,
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
		u.Error = err.Error()
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.URL, nil)
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err.Error()
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:58.0) Gecko/20100101 Firefox/58.0")

	resp, err := client.Do(req)
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err.Error()
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		w.logger.Println(resp.Status)
		return errors.New("Error fetching page")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err.Error()
		return err
	}

	u.FetchTime = time.Now()

	s := sha256.Sum256(body)
	sum := hex.EncodeToString(s[:])

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		w.logger.Println(err)
		w.instrument.Gauge("workers_active", -1)
		u.Error = err.Error()
		return err
	}

	title := doc.Find("title").First().Text()
	description, _ := doc.Find("meta[name=description]").First().Attr("content")

	urls := doc.Find("a[href]").Map(func(i int, sel *goquery.Selection) string {
		if attr, ok := sel.Attr("href"); ok {
			return attr
		}
		return ""
	})

	urls2 := doc.Find("area[href]").Map(func(i int, sel *goquery.Selection) string {
		if attr, ok := sel.Attr("href"); ok {
			return attr
		}
		return ""
	})

	urls = append(urls, urls2...)

	normalisedUrls := w.normaliseUrls(urls, u.URL)

	metadata := make(map[string]interface{})
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
			if name == "article:tag" {
				tagsIn, ok := metadata[name]
				if !ok {
					metadata[name] = []string{value}
				} else {
					tags := tagsIn.([]string)
					metadata[name] = append(tags, value)
				}
			} else {
				metadata[name] = value
			}
		}
	})
	doc.Find("meta[property]").Each(func(_ int, sel *goquery.Selection) {
		var name string
		var value string

		if attr, ok := sel.Attr("property"); ok {
			name = attr
		}

		if attr, ok := sel.Attr("content"); ok {
			value = strings.TrimSpace(attr)
		}

		if name != "" && value != "" {
			if name == "article:tag" {
				tagsIn, ok := metadata[name]
				if !ok {
					metadata[name] = []string{strings.TrimSpace(value)}
				} else {
					tags := tagsIn.([]string)
					metadata[name] = append(tags, strings.TrimSpace(value))
				}

			} else {
				metadata[name] = value
			}
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
		u.Error = err.Error()
	}

	var harvested []interface{}
	fncs := w.extractors.Matches(parsedURL.Host)
	for _, fn := range fncs {
		h, err := fn.Extract(doc)
		if err != nil {
			w.logger.Println(err)
			continue
		}

		if h != nil {
			harvested = append(harvested, h)
		}
	}

	u.ExtractTime = time.Now()
	u.PageHash = sum

	u.Title = title
	u.Description = description
	u.MicroData = mdata
	u.HarvestedData = harvested
	u.HarvestedURLs = normalisedUrls
	u.JSONData = jld
	u.MetaData = metadata

	w.instrument.Gauge("workers_active", -1)
	w.instrument.Count("crawl_url")

	u.EndTime = time.Now()

	w.publisher.Publish(u)
	return nil
}

func (w *defaultWorker) normaliseUrls(urls []string, ref string) []string {
	out := []string{}
	for _, u := range urls {
		retURL, err := url.Parse(u)
		if err != nil {
			continue
		}

		refURL, err := url.Parse(ref)
		if err != nil {
			continue
		}

		nu := refURL.ResolveReference(retURL)

		normalised, err := purell.NormalizeURLString(nu.String(), purell.FlagsSafe|purell.FlagRemoveFragment)
		if err != nil {
			continue
		}

		out = append(out, normalised)
	}

	//dedup
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range out {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	//remove javascript urls
	list2 := []string{}
	for _, entry := range out {
		if !strings.Contains(entry, "javascript:") {
			list2 = append(list2, entry)
		}
	}

	return list2
}
