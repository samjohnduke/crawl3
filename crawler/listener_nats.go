package crawler

import (
	"encoding/json"
	"log"

	nats "github.com/nats-io/go-nats"
)

// ListenerNats implements the Listener interface over the nats message bus
type ListenerNats struct {
	nc  *nats.Conn
	c   chan *Crawl
	sub *nats.Subscription
}

// NewListenerNats creates a new Nats listener for observing the output
// of successful crawls
func NewListenerNats(nc *nats.Conn) Listener {
	return &ListenerNats{
		nc:  nc,
		c:   make(chan *Crawl, 100),
		sub: nil,
	}
}

// Listen creates a channel that you can use to consume crawls
func (ln *ListenerNats) Listen() chan *Crawl {
	var err error
	ln.sub, err = ln.nc.Subscribe("crawl_complete", func(msg *nats.Msg) {
		var c *Crawl
		err := json.Unmarshal(msg.Data, &c)
		if err != nil {
			log.Println(err)
			return
		}

		ln.c <- c
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	return ln.c
}

// Close the nats connection
func (ln *ListenerNats) Close() {

}
