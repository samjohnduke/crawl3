package aggregator

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/iancoleman/strcase"

	driver "github.com/arangodb/go-driver"
	"github.com/davecgh/go-spew/spew"
	"github.com/samjohnduke/crawl3/crawler"
)

// Service implementations provide a view over the entire dataset of the application
type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
	Query(ctx context.Context, queryStr string) (result interface{}, err error)
}

// Client implementations provide a way for others to read the data
type Client interface {
	Query(ctx context.Context, queryStr string) (result interface{}, err error)
}

// Query is a data object to hold the query string
type Query struct {
	query string
}

// Aggregator is responsible for outputing the data requested by a client
type Aggregator struct {
	in            chan *crawler.Crawl
	aql           driver.Client
	db            driver.Database
	meta          driver.Collection
	content       driver.Collection
	pages         driver.Collection
	entities      driver.Collection
	entityRef     driver.Collection
	keyRef        driver.Collection
	namedKeywords driver.Collection
	uriRef        driver.Collection
	metaRef       driver.Collection
	quit          chan chan bool
}

// Opts allows you to pass configuration options into the aggregator
type Opts struct {
	Listener     chan *crawler.Crawl
	ArangoClient driver.Client
}

// New creates a new aggregator from the provided options
func New(opts Opts) (*Aggregator, error) {
	db, err := opts.ArangoClient.Database(context.Background(), "crawl3")
	if err != nil {
		return nil, err
	}

	metaCol, err := db.Collection(context.Background(), "meta")
	if err != nil {
		return nil, err
	}

	pages, err := db.Collection(context.Background(), "pages")
	if err != nil {
		return nil, err
	}

	content, err := db.Collection(context.Background(), "content")
	if err != nil {
		return nil, err
	}

	entity, err := db.Collection(context.Background(), "entities")
	if err != nil {
		return nil, err
	}

	entityRef, err := db.Collection(context.Background(), "entity_ref")
	if err != nil {
		return nil, err
	}

	keyRef, err := db.Collection(context.Background(), "key_ref")
	if err != nil {
		return nil, err
	}

	nk, err := db.Collection(context.Background(), "named_keywords")
	if err != nil {
		return nil, err
	}

	uriRef, err := db.Collection(context.Background(), "uri_ref")
	if err != nil {
		return nil, err
	}

	metaRef, err := db.Collection(context.Background(), "meta_ref")
	if err != nil {
		return nil, err
	}

	return &Aggregator{
		in:            opts.Listener,
		aql:           opts.ArangoClient,
		db:            db,
		meta:          metaCol,
		content:       content,
		pages:         pages,
		entities:      entity,
		uriRef:        uriRef,
		keyRef:        keyRef,
		namedKeywords: nk,
		entityRef:     entityRef,
		metaRef:       metaRef,
		quit:          make(chan chan bool),
	}, nil
}

// Query returns some data based of a query (query lang yet to be decided)
func (a *Aggregator) Query(ctx *context.Context, q string) (interface{}, error) {
	return nil, nil
}

// Start listening to the
func (a *Aggregator) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case c := <-a.in:
				a.store(c)
			case q := <-a.quit:
				q <- true
				return
			}
		}
	}()

	return nil
}

// Stop the aggregator
func (a *Aggregator) Stop(ctx context.Context) error {
	wait := make(chan bool)
	a.quit <- wait

	select {
	case <-wait:
		return nil
	case <-ctx.Done():
		return errors.New("Unable to stop the aggregator in time")
	}
}

func (a *Aggregator) store(c *crawler.Crawl) {
	var metaDescription string
	var metaTitle string

	if d, exists := c.MetaData["description"]; exists {
		metaDescription = d
	}

	if t, exists := c.MetaData["title"]; exists {
		metaTitle = t
	}

	doc := document{
		URL:         c.URL,
		Title:       metaTitle,
		Description: metaDescription,
		Data:        c.HarvestedData,
		JData:       c.JSONData,
		MData:       c.MicroData,
	}

	ref, ok := a.documentExists(c.URL)
	if !ok {
		m, err := a.pages.CreateDocument(context.Background(), doc)
		if err != nil {
			log.Println(err)
			return
		}

		ref = &m
	} else {
		m, err := a.pages.UpdateDocument(context.Background(), ref.Key, doc)
		if err != nil {
			log.Println(err)
			return
		}

		ref = &m
	}

	a.processMeta(c, ref.ID.String())
	a.processLinks(c, ref.ID.String())
}

func (a *Aggregator) processMeta(c *crawler.Crawl, ref string) {
	var metaCards = []map[string]string{}

	for _, tag := range keywordMetaTags {
		if contents, exists := c.MetaData[tag]; exists {
			ks := strings.Split(contents, ",")

			for _, k := range ks {
				a.addKeyWord(k, ref)
			}
		}
	}

	for _, tag := range colonCards {
		card := make(map[string]string)
		for k, v := range c.MetaData {
			if strings.HasPrefix(k, tag+":") {
				card[strings.TrimPrefix(k, tag+":")] = v
			}
		}
		card["type"] = tag
		metaCards = append(metaCards, card)
	}
	for _, tag := range dotCards {
		card := make(map[string]string)
		for k, v := range c.MetaData {
			if strings.HasPrefix(k, tag+".") {
				card[strings.TrimPrefix(k, tag+".")] = v
			}
		}
		card["type"] = tag
		metaCards = append(metaCards, card)
	}

	for _, card := range metaCards {
		a.addCard(card, ref)
	}
}

func (a *Aggregator) addKeyWord(key string, ref string) error {
	key = strings.TrimSpace(key)
	key = strings.ToLower(key)
	keyref := strcase.ToSnake(key)

	kw := keyword{
		Keyword: key,
		Key:     keyref,
	}

	exists, err := a.namedKeywords.DocumentExists(context.Background(), keyref)
	if err != nil {
		log.Println(err)
	}

	var m driver.DocumentMeta
	if !exists {
		m, err = a.namedKeywords.CreateDocument(context.Background(), kw)
		if err != nil {
			log.Println(err)
		}
	} else {
		var k keyword
		m, err = a.namedKeywords.ReadDocument(context.Background(), keyref, &k)
		if err != nil {
			log.Println(err)
		}
	}

	_, err = a.keyRef.CreateDocument(context.Background(), map[string]string{
		"_to":   m.ID.String(),
		"_from": ref,
		"name":  "has_keyword",
		"date":  time.Now().Format(time.RFC3339),
	})

	if err != nil {
		spew.Dump(err)
		return err
	}

	return nil
}

func (a *Aggregator) addCard(card map[string]string, ref string) error {
	m, err := a.meta.CreateDocument(context.Background(), card)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = a.metaRef.CreateDocument(context.Background(), map[string]string{
		"_from": ref,
		"_to":   m.ID.String(),
		"type":  card["type"],
	})

	return err
}

func (a *Aggregator) processLinks(c *crawler.Crawl, ref string) {
	for _, u := range c.HarvestedURLs {

		m, ok := a.documentExists(u)
		if !ok {
			me, err := a.pages.CreateDocument(context.Background(), map[string]interface{}{
				"URL": u,
			})
			if err != nil {
				log.Println(err)
				continue
			}
			m = &me
		}

		_, err := a.uriRef.CreateDocument(context.Background(), map[string]string{
			"_from": ref,
			"_to":   m.ID.String(),
			"type":  "links_to",
		})

		if err != nil {
			log.Println(err)
			continue
		}
	}
}

var keywordMetaTags = []string{"keywords", "news_keywords"}
var colonCards = []string{"og", "twitter", "article", "fb"}
var dotCards = []string{"ABC", "DCTERMS", "geo"}

type keyword struct {
	Keyword string `json:"keyword"`
	Key     string `json:"_key"`
}

type keyRef struct {
	To   string `json:"_to"`
	From string `json:"_from"`
}

type entity struct {
}

type content struct {
}

type document struct {
	URL         string `json:"URL"`
	Title       string
	Description string
	Data        interface{}
	JData       interface{}
	MData       interface{}
}

func (a *Aggregator) documentExists(k string) (*driver.DocumentMeta, bool) {
	query := "FOR d IN pages FILTER d.URL == @url LIMIT 1 RETURN d"
	cursor, err := a.db.Query(context.Background(), query, map[string]interface{}{
		"url": k,
	})
	if err != nil {
		return nil, false
	}
	defer cursor.Close()

	for {
		var doc document
		meta, err := cursor.ReadDocument(context.Background(), &doc)
		if driver.IsNoMoreDocuments(err) {
			return nil, false
		} else if err != nil {
			return nil, false
		}

		return &meta, true
	}

}
