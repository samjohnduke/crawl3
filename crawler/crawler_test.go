package crawler

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestNewService(t *testing.T) {
	service, err := New(ServiceOpts{}, func(opts WorkerOpts) WorkerFactoryFunc {
		return func(pool chan chan *Crawl) Worker {
			return NewDefaultWorker(pool, opts)
		}
	})
	if err != nil {
		t.Error(err)
	}

	result, err := service.Crawl(context.Background(), "https://google.com")
	if err != nil {
		t.Error(err)
	}

	spew.Dump(result)
}
