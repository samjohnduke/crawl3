package schedular

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/samjohnduke/crawl3/shared"

	"github.com/samjohnduke/crawl3/crawler"
)

// SitemapSchedular manages the scheduling of abc pages
type SitemapSchedular struct {
	sched           Service
	store           Store
	quit            chan chan bool
	sitemaps        []string
	excludeSitemaps map[string]struct{}
	sitemapFilter   SitemapFilter
	lastRun         time.Time
}

// SitemapFilter is a function used to determine if the
type SitemapFilter func(string) bool

// SitemapSchedularOpts is used to configure the SitemapSchedular
type SitemapSchedularOpts struct {
	Sitemaps        []string
	ExcludeSitemaps map[string]struct{}
	SitemapFilter   SitemapFilter
}

func NewSitemapSchedularOpts(opts shared.HostSchedularOpts) SitemapSchedularOpts {
	var sitemaps []string
	for _, f := range opts.Data["sitemaps"].([]interface{}) {
		sitemaps = append(sitemaps, f.(string))
	}

	var filterFunc SitemapFilter
	ai, ok := opts.Data["filter"].([]interface{})
	if ok {
		for _, filter := range ai {
			ff, ok := filter.(map[string]interface{})
			if ok {
				for k, v := range ff {
					switch k {
					case "contains":
						filterFunc = func(u string) bool {
							arr := v.([]interface{})
							var r bool

							for _, s := range arr {
								if strings.Contains(u, s.(string)) {
									r = true
								}
							}

							return r
						}
						break
					}
				}
			}
		}
	}

	return SitemapSchedularOpts{
		Sitemaps:      sitemaps,
		SitemapFilter: filterFunc,
	}
}

// NewSitemapSchedular creates a new SitemapSchedular from the provided options
func NewSitemapSchedular(opts SitemapSchedularOpts) *SitemapSchedular {
	return &SitemapSchedular{
		sitemaps:        opts.Sitemaps,
		excludeSitemaps: opts.ExcludeSitemaps,
		sitemapFilter:   opts.SitemapFilter,
		quit:            make(chan chan bool),
	}
}

// Start takes a schedular and a store and begins loading urls into the schedular
func (s *SitemapSchedular) Start(sched Service, store Store) error {
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
				close(closer)
				return
			}
		}
	}()
	return nil
}

func (s *SitemapSchedular) run() {
	for _, sm := range s.sitemaps {
		sitemap, err := ParseSitemapFromURL(sm)
		if err != nil {
			log.Println(err)
			return
		}

		for _, in := range sitemap.SitemapIndex.Sitemap {
			if (time.Time{}).Equal(s.lastRun) {
				s.lastRun = time.Now().Add(-24 * time.Hour)
			}

			if !s.lastRun.Before(in.Lastmod.Time) {
				continue
			}

			_, ok := s.excludeSitemaps[in.Loc]
			if ok && s.excludeSitemaps != nil {
				continue
			}

			if s.sitemapFilter != nil && !s.sitemapFilter(in.Loc) {
				continue
			}

			resp, err := http.Get(in.Loc)
			if err != nil {
				log.Println(err)
				continue
			}
			defer resp.Body.Close()

			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(in.Loc)
				continue
			}

			contentType := http.DetectContentType(out)

			if contentType == "application/x-gzip" {
				var gr io.ReadCloser
				gr, err = gzip.NewReader(bytes.NewBuffer(out))
				defer gr.Close()
				out, err = ioutil.ReadAll(gr)
			}

			mm, err := ParseSitemap(out)
			if err != nil {
				log.Println(in.Loc)
				continue
			}

			sitemap.URLset.URLS = append(sitemap.URLset.URLS, mm.URLset.URLS...)
		}

		for _, u := range sitemap.URLset.URLS {
			if u.Lastmod.Time.Before(s.lastRun) {
				continue
			}

			if !s.store.HasVisited(u.Loc) && !s.store.IsQueued(u.Loc) {
				s.store.Queue(u.Loc)
				s.sched.Schedule(u.Loc)
			}
		}
	}

	s.lastRun = time.Now()
}

// Stop will turn off the schedular when we are done
func (s *SitemapSchedular) Stop(ctx context.Context) error {
	closer := make(chan bool)
	s.quit <- closer
	<-closer
	return nil
}

// Schedule recieves urls from the main schedular after a crawl is complete
func (s *SitemapSchedular) Schedule(c *crawler.Crawl) {}
