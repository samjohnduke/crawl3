package schedular

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/samjohnduke/crawl3/shared"

	"github.com/samjohnduke/crawl3/crawler"
)

// ErrCancelSchedule - Return this error if you want to manage the scheduling of harvested data
var ErrCancelSchedule = errors.New("Cancel Schedule for harvest")

// Service is the interface for launching the spider against a particular service border
// and the implmentation should use this appropriately
type Service interface {
	Start()
	Stop(context.Context) error
	Schedule(rootURL string) error
	ScheduleAfter(t time.Time, rootURL string) error
	OnHarvest(func(*crawler.Crawl) error)
}

// The Schedular is responsible for managing the application state of a crawl. It should
// listen for crawls, pull out the harvested urls, check for whether or not it is needed
// to crawl and then if yes intelligently push them to the crawler
type Schedular struct {
	visited    map[string]*shared.URL
	pending    map[string]*shared.URLList
	cb         func(*crawler.Crawl) error
	delay      time.Duration
	die        chan chan bool
	instrument crawler.Instrument
	logger     *log.Logger
	client     crawler.Client
	allowed    map[string]bool
}

// Opts are used to customise the
type Opts struct {
	Instrument     crawler.Instrument
	Logger         *log.Logger
	Client         crawler.Client
	CrawlDelay     time.Duration
	AllowedDomains []string
}

// A HostSchedular looks after a spefic host and is control of scheduling that hosts
// urls. This enables us to not just crawl but do via other means, such as
type HostSchedular interface {
	Start(Service, Store) error
	Stop(context.Context) error
	Schedule(*crawler.Crawl)
}

// NewHostSchedular builds the host schedular from the provided options,
// and using the opts.Type field, create the schedular of the correct type
func NewHostSchedular(opts shared.HostSchedularOpts) HostSchedular {
	switch opts.Type {
	case shared.RSS:
		return NewRSSSchedular(NewRSSSchedularOpts(opts))
	case shared.Sitemap:
		return NewSitemapSchedular(NewSitemapSchedularOpts(opts))
	}
	return nil
}

// NewSchedular created the new schedular from a set of options
func NewSchedular(opts Opts) (Service, error) {
	allowed := make(map[string]bool)
	for _, a := range opts.AllowedDomains {
		allowed[a] = true
	}

	return &Schedular{
		visited:    make(map[string]*shared.URL),
		pending:    make(map[string]*shared.URLList),
		logger:     opts.Logger,
		instrument: opts.Instrument,
		client:     opts.Client,
		delay:      opts.CrawlDelay,
		die:        make(chan chan bool),
		allowed:    allowed,
	}, nil
}

// Start the schedular
func (s *Schedular) Start() {
	go func() {
		timer := time.NewTicker(s.delay)

		for {
			select {
			case <-timer.C:
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

	toCrawl := []*shared.URL{}
	for _, list := range s.pending {
		if list.Len() > 0 {
			toCrawl = append(toCrawl, list.Pop(1)...)
		}
	}

	for _, u := range toCrawl {
		if _, exists := s.visited[u.Normalised()]; exists {
			continue
		}

		rHost := u.Hostname()
		if _, exists := s.allowed[rHost]; exists {
			go s.crawl(u)
		}
	}
}

func (s *Schedular) crawl(u *shared.URL) {
	s.instrument.Gauge("scheduled_crawl_progress", 1)
	defer s.instrument.Gauge("scheduled_crawl_progress", -1)

	log.Println("starting crawl of: ", u.Normalised())
	crawl, err := s.client.Crawl(context.Background(), u.Normalised())
	if err != nil {
		s.instrument.Count("scheduled_crawl_error")
		log.Println(err)
		return
	}

	if crawl.Error != "" {
		log.Println(crawl.Error)
		return
	}

	s.visited[u.Normalised()] = u

	if s.cb != nil {
		err := s.cb(crawl)
		if err == ErrCancelSchedule {
			return
		}
	}

	for _, hu := range crawl.HarvestedURLs {
		url, err := shared.NewURLWithReference(hu, crawl.URL)
		if err != nil {
			log.Println(err)
			continue
		}
		s.Schedule(url.Normalised())
	}
}

// ScheduleAfter schedules a url to run after a given point in time
func (s *Schedular) ScheduleAfter(t time.Time, rootURL string) error {
	return nil
}

// Schedule pushes just a single url into the application
func (s *Schedular) Schedule(root string) error {
	rURL, err := shared.NewURL(root)
	if err != nil {
		return err
	}

	rHost := rURL.Hostname()

	s.instrument.Count("schedular_one")
	s.instrument.Histogram("schedular_one_host", rHost)

	a, exists := s.pending[rHost]
	if !exists {
		li := shared.NewURLList()
		li.Unshift(rURL)
		s.pending[rHost] = li
	} else {
		a.Unshift(rURL)
	}

	return nil
}

// OnHarvest is called when data has been harvest.
func (s *Schedular) OnHarvest(fn func(*crawler.Crawl) error) {
	s.cb = fn
}
