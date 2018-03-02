package crawler

import (
	"context"
	"log"
	"testing"
	"time"

	nats "github.com/nats-io/go-nats"
)

func TestCrawlClient(t *testing.T) {
	ser := RunDefaultServer()
	defer ser.Shutdown()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
	}

	service, err := New(ServiceOpts{}, func(opts WorkerOpts) WorkerFactoryFunc {
		return func(pool chan chan *Crawl) Worker {
			return NewDefaultWorker(pool, opts)
		}
	})
	if err != nil {
		t.Error(err)
	}

	transport := NewTransportNats(nc, service)

	err = transport.Start(context.Background())
	if err != nil {
		t.Error(err)
	}

	client := NewClientNats(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = client.Crawl(ctx, "https://google.com")
	if err != nil {
		log.Println(err)
		t.Error(err)

		err = transport.Stop(context.Background())
		if err != nil {
			t.Error(err)
		}

		return
	}

	err = transport.Stop(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestCrawlAsyncClient(t *testing.T) {
	ser := RunDefaultServer()
	defer ser.Shutdown()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
	}

	service, err := New(ServiceOpts{}, func(opts WorkerOpts) WorkerFactoryFunc {
		return func(pool chan chan *Crawl) Worker {
			return NewDefaultWorker(pool, opts)
		}
	})
	if err != nil {
		t.Error(err)
	}

	transport := NewTransportNats(nc, service)

	err = transport.Start(context.Background())
	if err != nil {
		t.Error(err)
	}

	client := NewClientNats(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, test := range testUrls {
		wait := make(chan bool)
		_, err = client.CrawlAsync(ctx, test, func(c *Crawl) {
			wait <- true
		})
		if err != nil {
			log.Println(err)
			t.Error(err)

			err = transport.Stop(context.Background())
			if err != nil {
				t.Error(err)
			}

			return
		}

		<-wait
	}

	err = transport.Stop(context.Background())
	if err != nil {
		t.Error(err)
	}
}

var testUrls = []string{
	"https://google.com",
	"https://facebook.com",
}
