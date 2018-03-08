package aggregator

import (
	"context"
	"errors"

	driver "github.com/arangodb/go-driver"
	"github.com/samjohnduke/crawl3/crawler"
)

// Service implementations provide a view over the entire dataset of the application
type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
	Query(ctx context.Context, queryStr string) (result interface{}, err error)
}

// Client implementations provide a way for others to read the data
type Client interface {
	Query(ctx context.Context, queryStr string) (result interface{}, err error)
}

// Query is a data object to hold the query string
type Query struct {
	query string
}

// Aggregator is responsible for outputing the data requested by a client
type Aggregator struct {
	in   chan *crawler.Crawl
	aql  driver.Client
	quit chan chan bool
}

// Opts allows you to pass configuration options into the aggregator
type Opts struct {
	Listener     chan *crawler.Crawl
	ArangoClient driver.Client
}

// New creates a new aggregator from the provided options
func New(opts Opts) *Aggregator {
	return &Aggregator{
		in:   opts.Listener,
		aql:  opts.ArangoClient,
		quit: make(chan chan bool),
	}
}

// Query returns some data based of a query (query lang yet to be decided)
func (a *Aggregator) Query(ctx *context.Context, q string) (interface{}, error) {
	return nil, nil
}

// Start listening to the
func (a *Aggregator) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case c := <-a.in:
				a.store(c)
			}
		}
	}()

	return nil
}

func (a *Aggregator) Stop(ctx context.Context) error {
	wait := make(chan bool)
	a.quit <- wait

	select {
	case <-wait:
		return nil
	case <-ctx.Done():
		return errors.New("Unable to stop the aggregator in time")
	}

	return nil
}

func (a *Aggregator) store(c *crawler.Crawl) {

}
