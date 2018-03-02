package crawler

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/gnatsd/server"
	nats "github.com/nats-io/go-nats"
)

// DefaultTestOptions are default options for the unit tests.
var DefaultTestOptions = server.Options{
	Host:           "localhost",
	Port:           4223,
	NoLog:          true,
	NoSigs:         true,
	MaxControlLine: 256,
}

// RunDefaultServer starts a new Go routine based server using the default options
func RunDefaultServer() *server.Server {
	return RunServer(&DefaultTestOptions)
}

// RunServer starts a new Go routine based server
func RunServer(opts *server.Options) *server.Server {
	if opts == nil {
		opts = &DefaultTestOptions
	}
	s := server.New(opts)
	if s == nil {
		panic("No NATS Server object returned.")
	}

	// Run server in Go routine.
	go s.Start()

	// Wait for accept loop(s) to be started
	if !s.ReadyForConnections(10 * time.Second) {
		panic("Unable to start NATS Server in Go Routine")
	}
	return s
}

func SetupTransport(nc *nats.Conn) (chan bool, error) {
	service, err := New(ServiceOpts{}, func(opts WorkerOpts) WorkerFactoryFunc {
		return func(pool chan chan *Crawl) Worker {
			return NewDefaultWorker(pool, opts)
		}
	})
	if err != nil {
		return nil, err
	}

	transport := NewTransportNats(nc, service)

	err = transport.Start(context.Background())
	if err != nil {
		return nil, err
	}

	close := make(chan bool)

	go func() {
		<-close
		transport.Stop(context.Background())
	}()

	return close, nil
}

func TearDownTransport() {

}

func TestNatsTransportSetup(t *testing.T) {
	ser := RunDefaultServer()
	defer ser.Shutdown()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
	}

	close, err := SetupTransport(nc)
	if err != nil {
		t.Error(err)
	}

	close <- true
}

func TestNatsTransportCrawl(t *testing.T) {
	ser := RunDefaultServer()
	defer ser.Shutdown()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Error(err)
	}

	close, err := SetupTransport(nc)
	if err != nil {
		t.Error(err)
	}

	req := CrawlRequest{
		URL: "https://google.com",
	}

	out, _ := json.Marshal(req)

	_, err = nc.Request("crawl", out, 10*time.Second)
	if err != nil {
		t.Error(err)
	}

	close <- true
}
