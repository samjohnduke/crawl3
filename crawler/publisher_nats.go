package crawler

import (
	"encoding/json"

	nats "github.com/nats-io/go-nats"
)

type publisherNats struct {
	conn *nats.Conn
}

// NewPublisherNats creates the transport for the nats service
func NewPublisherNats(conn *nats.Conn) Publisher {
	return &publisherNats{
		conn: conn,
	}
}

// Publish pushes a completed crawl onto the message bus
func (p *publisherNats) Publish(c *Crawl) error {
	d, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = p.conn.Publish("crawl_complete", d)
	if err != nil {
		return err
	}

	return nil
}
