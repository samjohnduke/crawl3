package aggregator

import (
	"context"
	"encoding/json"

	nats "github.com/nats-io/go-nats"
	"github.com/samjohnduke/crawl3/crawler"
)

// TransportNats allows users to connect to the aggregator via
type TransportNats struct {
	conn    *nats.Conn
	service Service
}

// TransportNats is the implementation of the crawl transport via the nats
// message bus.
type transportNats struct {
	conn    *nats.Conn
	service Service
	subs    []*nats.Subscription
}

// NewTransportNats creates the transport for the nats service
func NewTransportNats(conn *nats.Conn, service Service) crawler.Transport {
	return &transportNats{
		conn: conn, service: service,
	}
}

// Start subscribes to the nats channels for the service
func (t *transportNats) Start(ctx context.Context) error {
	sub, err := t.conn.Subscribe("crawl", func(msg *nats.Msg) {
		var aq *Query
		err := json.Unmarshal(msg.Data, aq)
		if err != nil {
			return
		}

		t.service.Query(context.Background(), aq.query)
	})
	if err != nil {
		return err
	}
	t.subs = append(t.subs, sub)

	return nil
}

// Stop closes all the subscriptions
func (t *transportNats) Stop(ctx context.Context) error {
	for _, sub := range t.subs {
		err := sub.Unsubscribe()
		if err != nil {
			return err
		}
	}
	return nil
}
