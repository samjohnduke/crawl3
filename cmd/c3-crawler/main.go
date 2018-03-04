package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
)

func main() {
	log.Println("Setting up Crawler")

	//Setup the system to wait for shutdown
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//Connect to the nats server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	instrument := crawler.NewInstrumentationMem()
	go func() {
		spew.Dump(instrument)
	}()

	//create the service
	service, err := crawler.New(crawler.ServiceOpts{
		Instrument: instrument,
		Logger:     log.New(os.Stdout, "", log.LstdFlags),
	}, func(opts crawler.WorkerOpts) crawler.WorkerFactoryFunc {
		return func(pool chan chan *crawler.Crawl) crawler.Worker {
			return crawler.NewDefaultWorker(pool, opts)
		}
	})
	if err != nil {
		log.Fatal(err)
	}

	//Connect the service to the nats transport
	transport := crawler.NewTransportNats(nc, service)

	err = transport.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("crawler started")

	// Wait for shutdown and then turn off the transport
	go func() {
		<-sigs
		log.Println("Shutting Down Crawler")

		err = transport.Stop(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	<-done
}
