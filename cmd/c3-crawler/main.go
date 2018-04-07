package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
	"github.com/samjohnduke/crawl3/shared"
)

var setupMsg = "Setting up Crawler"
var stopMsg = "Stopping Crawler"
var workerCount int
var hostDir string

func main() {
	log.Println(setupMsg)

	flag.StringVar(&hostDir, "hostDir", "../models", "The directory that stores the models")
	flag.IntVar(&workerCount, "wc", 40, "The number of workers to spin up")
	flag.Parse()

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
	publisher := crawler.NewPublisherNats(nc)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	execs := crawler.NewDefaultExtractors()

	// load models for extracting data
	hosts, err := shared.LoadHostsFromDir(hostDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range hosts {
		e := &crawler.JSONExtractor{Rules: host.Extractor}
		execs.Add(e)
	}

	//create the service
	service, err := crawler.New(crawler.ServiceOpts{
		Instrument:  instrument,
		Logger:      logger,
		WorkerCount: int64(workerCount),
		Publisher:   publisher,
		Extractors:  execs,
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

	// Wait for shutdown and then turn off the transport
	go func() {
		<-sigs
		log.Println(stopMsg)

		err = transport.Stop(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	<-done
}
