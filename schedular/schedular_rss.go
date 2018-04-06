package schedular

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/samjohnduke/crawl3/shared"

	"github.com/mmcdole/gofeed"
	"github.com/samjohnduke/crawl3/crawler"
)

// RSSSchedular manages the scheduling of abc pages
type RSSSchedular struct {
	sched         Service
	store         Store
	quit          chan chan bool
	feeds         []string
	allowInsecure bool
}

type RSSSchedularOpts struct {
	Feeds         []string
	AllowInsecure bool
}

func NewRSSSchedularOpts(opts shared.HostSchedularOpts) RSSSchedularOpts {
	var feeds []string
	for _, f := range opts.Data["feeds"].([]interface{}) {
		feeds = append(feeds, f.(string))
	}

	var insecure bool
	ai, ok := opts.Data["allow_insecure"]
	if ok {
		insecure = ai.(bool)
	}

	return RSSSchedularOpts{
		Feeds:         feeds,
		AllowInsecure: insecure,
	}
}

func NewRSSSchedular(opts RSSSchedularOpts) *RSSSchedular {
	return &RSSSchedular{
		feeds:         opts.Feeds,
		allowInsecure: opts.AllowInsecure,
		quit:          make(chan chan bool),
	}
}

// Start takes a schedular and a store and begins loading urls into the schedular
func (s *RSSSchedular) Start(sched Service, store Store) error {
	s.sched = sched
	s.store = store

	go func() {
		go s.run()
		ticker := time.NewTicker(15 * time.Minute)
		for {
			select {
			case <-ticker.C:
				s.run()
				break

			case closer := <-s.quit:
				closer <- true
				close(closer)
				return
			}
		}
	}()
	return nil
}

func (s *RSSSchedular) run() {
	for _, f := range s.feeds {
		fp := gofeed.NewParser()

		tr := &http.Transport{}
		if s.allowInsecure {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}

		client := &http.Client{Transport: tr}

		resp, err := client.Get(f)
		if err != nil {
			log.Println(err, f)
			continue
		}
		defer resp.Body.Close()

		feed, err := fp.Parse(resp.Body)
		if err != nil {
			log.Println(err, f)
			continue
		}

		for _, i := range feed.Items {
			var link string
			if _, ok := i.Extensions["feedburner"]; ok {
				link = i.Extensions["feedburner"]["origLink"][0].Value
			} else {
				link = i.Link
			}

			if !s.store.HasVisited(link) && !s.store.IsQueued(link) {
				s.store.Queue(link)
				s.sched.Schedule(link)
			}
		}
	}
}

// Stop will turn off the schedular when we are done
func (s *RSSSchedular) Stop(ctx context.Context) error {
	closer := make(chan bool)
	s.quit <- closer
	<-closer
	return nil
}

// Schedule recieves urls from the main schedular after a crawl is complete.
// The abc schedular doesn't accept schedule urls as we will load everything from rss
func (s *RSSSchedular) Schedule(c *crawler.Crawl) {
}
