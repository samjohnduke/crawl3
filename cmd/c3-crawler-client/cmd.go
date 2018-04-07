package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	nats "github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
)

func main() {
	var url string
	if len(os.Args) == 2 {
		url = os.Args[1]
	} else if len(os.Args) > 2 {
		log.Fatal("Too many arguments, please only pass a single url")
	} else {
		log.Fatal("missing url argument")
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	client := crawler.NewClientNats(nc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Crawl(ctx, url)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(result)
}
