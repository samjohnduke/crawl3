package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	nats "github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
	"github.com/samjohnduke/crawl3/schedular"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	client := crawler.NewClientNats(nc)

	opts := schedular.Opts{
		Client:               client,
		ConcurrencyPerDomain: 2,
		CrawlDelay:           1 * time.Second,
		Instrument:           crawler.NewInstrumentationMem(),
		Logger:               log.New(os.Stdout, log.Prefix(), log.LstdFlags),
		AllowedDomains:       []string{"www.paulgraham.com"},
	}

	sched, err := schedular.NewSchedular(opts)
	sched.Start()

	sched.OnHarvest(func(c *crawler.Crawl) error {
		return nil
	})

	sched.Schedule("http://www.paulgraham.com")

	//Setup the system to wait for shutdown
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting Down Schedular")

		err = sched.Stop(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	<-done
}
