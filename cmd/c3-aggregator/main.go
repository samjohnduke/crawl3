package main

import (
	"log"

	"github.com/samjohnduke/crawl3/aggregator"
	"github.com/samjohnduke/crawl3/crawler"

	nats "github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	l := crawler.NewListenerNats(nc)

	aggregator.New(aggregator.Opts{
		Listener: l.Listen(),
	})

}
