package main

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	nats "github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	client := crawler.NewClientNats(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Crawl(ctx, "https://google.com")
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(result)
}
