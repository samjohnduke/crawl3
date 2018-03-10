package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/samjohnduke/crawl3/aggregator"
	"github.com/samjohnduke/crawl3/crawler"

	nats "github.com/nats-io/go-nats"
)

func main() {
	log.Println("Starting Aggregator")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	l := crawler.NewListenerNats(nc)

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		log.Fatal(err)
	}
	aql, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("sam", "sammax1PaHd91U6"),
	})
	if err != nil {
		log.Fatal(err)
	}

	a, err := aggregator.New(aggregator.Opts{
		Listener:     l.Listen(),
		ArangoClient: aql,
	})

	if err != nil {
		log.Fatal(err)
	}

	err = a.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting Down Aggregator")

		l.Close()

		err = a.Stop(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		done <- true
	}()

	<-done
}
