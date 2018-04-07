package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samjohnduke/crawl3/shared"

	nats "github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
	"github.com/samjohnduke/crawl3/schedular"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var hostDir string
var sqldriver string
var sqlurl string
var startMsg = "Starting Scheduler"
var stopMsg = "Stopping Scheduler"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println(startMsg)

	flag.StringVar(&hostDir, "hostDir", "../models", "The directory that stores the models")
	flag.StringVar(&sqldriver, "sqldriver", "sqlite3", "The sql driver for storing data")
	flag.StringVar(&sqlurl, "sqlurl", "./dev.db", "the sql url use to connect to")
	flag.Parse()

	hosts, err := shared.LoadHostsFromDir(hostDir)
	if err != nil {
		log.Fatal(err)
	}

	var allowedHosts = []string{}
	var schedulers = make(map[string]schedular.HostSchedular)
	for _, host := range hosts {
		for _, o := range host.Schedular {
			schedulers[host.Host] = schedular.NewHostSchedular(o)
		}
		allowedHosts = append(allowedHosts, host.Host)
		allowedHosts = append(allowedHosts, host.Alias...)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	client := crawler.NewClientNats(nc)
	instrumentation := crawler.NewInstrumentationMem()
	crawlDelay := 2 * time.Second
	logger := log.New(os.Stdout, log.Prefix(), log.LstdFlags|log.Lshortfile)

	opts := schedular.Opts{
		Client:         client,
		CrawlDelay:     crawlDelay,
		Instrument:     instrumentation,
		Logger:         logger,
		AllowedDomains: allowedHosts,
	}

	sched, err := schedular.NewSchedular(opts)
	sched.Start()

	db, err := gorm.Open(sqldriver, sqlurl)
	if err != nil {
		log.Fatal(err)
	}

	store, err := schedular.NewSQLGormStore(db)
	if err != nil {
		log.Fatal(err)
	}

	// Start the reschedular
	cancelRescheduler := make(chan bool)
	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for {
			select {
			case <-ticker.C:
				store.Reschedule(sched)
				break

			case <-cancelRescheduler:
				return
			}
		}
	}()

	//Start all the schedulers
	for _, s := range schedulers {
		err := s.Start(sched, store)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Manage the harvested data
	sched.OnHarvest(func(c *crawler.Crawl) error {
		_, err := store.Visit(c.URL, c.PageHash)

		if err != nil {
			log.Println(err)
			return schedular.ErrCancelSchedule
		}

		for range c.HarvestedURLs {
			if s, ok := schedulers[c.Host()]; ok {
				s.Schedule(c)
			}
		}

		return schedular.ErrCancelSchedule
	})

	//Setup the system to wait for shutdown
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println(stopMsg)

		db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = sched.Stop(ctx)
		if err != nil {
			log.Fatal(err)
		}

		close(cancelRescheduler)

		for _, s := range schedulers {
			err := s.Stop(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}

		done <- true
	}()

	<-done
}
