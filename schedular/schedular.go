package schedular

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/samjohnduke/crawl3/crawler"
)

// ErrCancelSchedule - Return this error if you want to manage the scheduling of harvested data
var ErrCancelSchedule = errors.New("Cancel Schedule for harvest")

// Service is the interface for launching the spider against a particular service border
// and the implmentation should use this appropriately
type Service interface {
	Schedule(rootURL string) error
	ScheduleAfter(t time.Time, rootURL string) error
	OnHarvest(func(*crawler.Crawl) error)
}

// The Schedular is responsible for managing the application state of a crawl. It should
// listen for crawls, pull out the harvested urls, check for whether or not it is needed
// to crawl and then if yes intelligently push them to the crawler
type Schedular struct {
	visited    map[string]*URL
	pending    map[string]*urlList
	later      *TimeMap
	cb         func(*crawler.Crawl) error
	frequency  int
	delay      time.Duration
	die        chan chan bool
	instrument crawler.Instrument
	logger     *log.Logger
	client     crawler.Client
}

// Opts are used to customise the
type Opts struct {
	Instrument           crawler.Instrument
	Logger               *log.Logger
	Client               crawler.Client
	ConcurrencyPerDomain int
	CrawlDelay           time.Duration
}

// NewSchedular created the new schedular from a set of options
func NewSchedular(opts Opts) (Service, error) {
	return &Schedular{
		visited:    make(map[string]*URL),
		pending:    make(map[string]*urlList),
		later:      &TimeMap{},
		frequency:  opts.ConcurrencyPerDomain,
		logger:     opts.Logger,
		instrument: opts.Instrument,
		client:     opts.Client,
		delay:      opts.CrawlDelay,
		die:        make(chan chan bool),
	}, nil
}

// Start the schedular
func (s *Schedular) Start() {
	go func() {
		timer := time.NewTicker(s.delay)

		for {
			select {
			case <-timer.C:
				log.Println("Scheduling URLs")
				s.run()
				break

			case close := <-s.die:
				close <- true
				return
			}
		}
	}()
}

// Stop the schedular
func (s *Schedular) Stop(ctx context.Context) error {
	cancel := make(chan bool)
	s.die <- cancel

	select {
	case <-ctx.Done():
		return errors.New("Stop timed out")

	case <-cancel:
		log.Println("Schedular stopped")
	}
	return nil
}

func (s *Schedular) run() {
	toCrawl := []*URL{}
	for _, list := range s.pending {
		if list.Len() > s.frequency {
			toCrawl = append(toCrawl, list.pop(s.frequency)...)
		}
	}

	for _, u := range toCrawl {
		go s.crawl(u)
	}
}

func (s *Schedular) crawl(u *URL) {
	s.instrument.Gauge("scheduled_crawl_progress", 1)
	defer s.instrument.Gauge("scheduled_crawl_progress", -1)

	crawl, err := s.client.Crawl(context.Background(), u.normalised)
	if err != nil {
		s.instrument.Count("scheduled_crawl_error")
		log.Println(err)
	}

	if s.cb != nil {
		err := s.cb(crawl)
		if err == ErrCancelSchedule {
			return
		}
	}

	for _, hu := range crawl.HarvestedURLs {
		url, err := NewURLWithReference(hu, crawl.URL)
		if err != nil {
			continue
		}
		s.Schedule(url.normalised)
	}

	spew.Dump(crawl.HarvestedURLs)
}

// ScheduleAfter schedules a url to run after a given point in time
func (s *Schedular) ScheduleAfter(t time.Time, rootURL string) error {
	return nil
}

// Schedule pushes just a single url into the application
func (s *Schedular) Schedule(root string) error {
	rURL, err := NewURL(root)
	if err != nil {
		return err
	}

	rHost := rURL.url.Hostname()

	s.instrument.Count("schedular_one")
	s.instrument.Histogram("schedular_one_host", rHost)

	a, exists := s.pending[rHost]
	if !exists {
		li := newURLList()
		li.unshift(rURL)
		s.pending[rHost] = li
	}
	a.unshift(rURL)

	return nil
}

// OnHarvest is called when data has been harvest.
func (s *Schedular) OnHarvest(fn func(*crawler.Crawl) error) {
	s.cb = fn
}
