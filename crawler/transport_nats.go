package crawler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/go-nats"
)

// TransportNats is the implementation of the crawl transport via the nats
// message bus.
type transportNats struct {
	conn    *nats.Conn
	service Service
	subs    []*nats.Subscription
}

// NewTransportNats creates the transport for the nats service
func NewTransportNats(conn *nats.Conn, service Service) Transport {
	return &transportNats{
		conn: conn, service: service,
	}
}

// Start subscribes to the nats channels for the service
func (t *transportNats) Start(ctx context.Context) error {
	sub, err := t.conn.QueueSubscribe("crawl", "crawl_worker", t.recieveCrawlRequest)
	if err != nil {
		return err
	}
	t.subs = append(t.subs, sub)

	sub2, err := t.conn.QueueSubscribe("crawlAsync", "crawl_worker", t.recieveCrawlAsyncRequest)
	if err != nil {
		return err
	}
	t.subs = append(t.subs, sub2)

	sub3, err := t.conn.QueueSubscribe("crawlProgress", "crawl_worker", t.recieveCrawlProgressRequest)
	if err != nil {
		return err
	}
	t.subs = append(t.subs, sub3)

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

// Send a crawl to the publish endpoint that others can listen into
func (t *transportNats) Publish(ctx context.Context, crawl *Crawl) error {
	return nil
}

// process a crawl request synchronous request
func (t *transportNats) recieveCrawlRequest(m *nats.Msg) {
	var crawlRequest CrawlRequest
	err := json.Unmarshal(m.Data, &crawlRequest)
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	result, err := t.service.Crawl(ctx, crawlRequest.URL)
	if err != nil {
		log.Println(err)
		return
	}

	reply := &CrawlReply{
		Crawl: *result,
		Error: err,
	}

	out, err := json.Marshal(reply)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.conn.Publish(m.Reply, out)
	if err != nil {
		log.Println(err)
		return
	}
}

// process a crawl request asynchronously
func (t *transportNats) recieveCrawlAsyncRequest(m *nats.Msg) {
	var crawlRequest CrawlAsyncRequest
	err := json.Unmarshal(m.Data, &crawlRequest)
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	guid, err := t.service.CrawlAsync(ctx, crawlRequest.URL, func(c *Crawl) {
		reply := &CrawlReply{
			Crawl: *c,
			Error: nil,
		}
		out, err := json.Marshal(reply)
		if err != nil {
			log.Println(err)
			return
		}

		t.conn.Publish(crawlRequest.Reply, out)
	})

	reply := &CrawlReply{
		Crawl: Crawl{ID: guid},
		Error: err,
	}

	out, err := json.Marshal(reply)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.conn.Publish(m.Reply, out)

	if err != nil {
		log.Println(err)
		return
	}
}

func (t *transportNats) recieveCrawlProgressRequest(m *nats.Msg) {
	var progressRequest ProgressRequest
	err := json.Unmarshal(m.Data, &progressRequest)
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	result, err := t.service.CrawlProgress(ctx, progressRequest.GUID)

	reply := &CrawlReply{
		Crawl: *result,
		Error: err,
	}

	out, err := json.Marshal(reply)
	if err != nil {
		log.Println(err)
		return
	}

	err = t.conn.Publish(m.Reply, out)

	if err != nil {
		log.Println(err)
		return
	}
}
