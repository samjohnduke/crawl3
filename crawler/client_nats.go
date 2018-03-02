package crawler

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
)

type clientNats struct {
	conn *nats.Conn
}

// NewClientNats creates a client for the service over the nats message bus
func NewClientNats(conn *nats.Conn) Client {
	return &clientNats{
		conn: conn,
	}
}

//An asynchronous request for crawling a webpage. The callback will be called when the crawl has been finished
func (c *clientNats) CrawlAsync(ctx context.Context, url string, cb func(crawl *Crawl)) (guid string, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	reqid := u.String()

	req := CrawlAsyncRequest{
		URL:   url,
		Reply: reqid,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	msg, err := c.conn.RequestWithContext(ctx, "crawlAsync", data)
	if err != nil {
		return "", nil
	}

	var reply CrawlReply
	err = json.Unmarshal(msg.Data, &reply)
	if err != nil {
		log.Println(err)
		return
	}

	if reply.Error != nil {
		log.Println(reply.Error)
		return "", reply.Error
	}

	sub, err := c.conn.Subscribe(reqid, func(msg *nats.Msg) {
		var reply CrawlReply
		err := json.Unmarshal(msg.Data, &reply)
		if err != nil {
			log.Println(err)
			return
		}

		if reply.Error != nil {
			log.Println(reply.Error)
			return
		}

		cb(&(reply.Crawl))
	})

	if err != nil {
		log.Println(err)
		return
	}

	sub.AutoUnsubscribe(1)

	return reply.Crawl.ID, err
}

//A synchronous request for crawling a webpage
func (c *clientNats) Crawl(ctx context.Context, url string) (result *Crawl, err error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	reqid := u.String()

	req := CrawlAsyncRequest{
		URL:   url,
		Reply: reqid,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg, err := c.conn.RequestWithContext(ctx, "crawl", data)
	if err != nil {
		return
	}

	var reply CrawlReply
	err = json.Unmarshal(msg.Data, &reply)
	if err != nil {
		log.Println(err)
		return
	}

	if reply.Error != nil {
		log.Println(reply.Error)
		return nil, nil
	}

	return &reply.Crawl, nil
}

// Get the progress of a crawl
func (c *clientNats) CrawlProgress(ctx context.Context, guid string) (result *Crawl, err error) {
	return nil, errors.New("Unimplemented")
}
