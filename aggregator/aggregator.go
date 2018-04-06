package aggregator

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/sbl/ner"

	driver "github.com/arangodb/go-driver"
	"github.com/samjohnduke/crawl3/crawler"
)

var ext *ner.Extractor

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
	pageToContent driver.Collection
	pageToEntity  driver.Collection
	entities      driver.Collection
	entityRef     driver.Collection
	keyRef        driver.Collection
	namedKeywords driver.Collection
	uriRef        driver.Collection
	metaRef       driver.Collection
	quit          chan chan bool
	ne            *ner.Extractor
}

type Aggregation struct {
	crawl       *crawler.Crawl
	document    document
	documentRef string
	entities    []tag
	article     NewsArticle
	cards       []map[string]interface{}
	keywords    []string
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

	pageContent, err := db.Collection(context.Background(), "page_to_content")
	if err != nil {
		return nil, err
	}

	pageEntity, err := db.Collection(context.Background(), "page_to_entity")
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

	entityRef, err := db.Collection(context.Background(), "entity_to_content")
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
		pageToContent: pageContent,
		pageToEntity:  pageEntity,
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
	var err error
	ext, err = ner.NewExtractor("/home/sam/.local/share/MITIE-models/english/ner_model.dat")
	if err != nil {
		return err
	}

	a.ne = ext

	go a.p()

	return nil
}

func (a *Aggregator) p() {
	for {
		select {
		case c := <-a.in:
			a.store(c)
			break

		case q := <-a.quit:
			ext.Free()
			q <- true
			return
		}
	}
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

	if _, ok := c.HarvestedData.([]interface{}); ok {

		doc := document{
			URL:         c.URL,
			Title:       c.Title,
			Description: c.Description,
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

		aggr := &Aggregation{
			crawl:       c,
			document:    doc,
			documentRef: ref.ID.String(),
		}

		a.processMeta(aggr)
		a.processData(aggr)

		a.saveAggregate(aggr)
	}
}

func (a *Aggregator) processMeta(aggr *Aggregation) {
	var metaCards = []map[string]interface{}{}

	for _, tag := range keywordMetaTags {
		if contents, exists := aggr.crawl.MetaData[tag]; exists {
			aggr.keywords = strings.Split(contents.(string), ",")
		}
	}

	for _, tag := range colonCards {
		card := make(map[string]interface{})

		for k, v := range aggr.crawl.MetaData {
			if strings.HasPrefix(k, tag+":") {
				if strings.TrimPrefix(k, tag+":") == "tag" {
					vv := v.([]interface{})
					tags := []string{}
					for i := range vv {
						//check for multiple in tags
						tt := strings.TrimSpace(vv[i].(string))
						tt2 := strings.Split(tt, ",")
						tags = append(tags, tt2...)
					}
					card["tags"] = tags
				} else {
					card[strings.TrimPrefix(k, tag+":")] = v.(string)
				}
			}
		}

		card["type"] = tag
		metaCards = append(metaCards, card)
	}

	for _, tag := range dotCards {
		card := make(map[string]interface{})
		for k, v := range aggr.crawl.MetaData {
			if strings.HasPrefix(k, tag+".") {
				card[strings.TrimPrefix(k, tag+".")] = v
			}
		}
		card["type"] = tag
		metaCards = append(metaCards, card)
	}

	aggr.cards = metaCards
}

func (a *Aggregator) processData(aggr *Aggregation) {
	data := aggr.crawl.HarvestedData.([]interface{})
	for i := range data {
		d, ok := data[i].(map[string]interface{})
		if !ok {
			log.Println("data is bad")
			continue
		}
		articles := d["articles"]

		for _, aa := range articles.([]interface{}) {
			aaa := aa.(map[string]interface{})

			newsArticle := NewsArticle{
				ArticleBody: aaa["content"].(string),
				Headline:    strings.Join(strings.Fields(strings.TrimSpace(aaa["title"].(string))), " "),
				Description: aggr.crawl.Description,
				PageRef:     aggr.documentRef,
				Keywords:    aggr.keywords,
			}

			for i := range newsArticle.Keywords {
				newsArticle.Keywords[i] = strings.TrimSpace(newsArticle.Keywords[i])
			}

			if p, ok := aaa["published_time"]; ok {
				ts := p.(string)
				tt, err := dateparse.ParseAny(ts)
				if err == nil {
					newsArticle.DatePublished = tt
				}
			}

			if imgInt, ok := aaa["images"]; ok {
				imgSlice := []string{}
				if imgs, ok := imgInt.([]interface{}); ok {
					for _, i := range imgs {
						imgSlice = append(imgSlice, i.(string))
					}
				}
				newsArticle.Images = imgSlice
			}

			for _, card := range aggr.cards {
				if acard, ok := card["type"]; ok && acard == "article" {
					for tagName, tagVal := range card {
						if tagName == "published_time" {
							t, err := dateparse.ParseAny(tagVal.(string))
							if err != nil {
								continue
							}
							newsArticle.DatePublished = t
						}

						if tagName == "modified_time" {
							t, err := dateparse.ParseAny(tagVal.(string))
							if err != nil {
								continue
							}
							newsArticle.DateModified = t
						}

						if tagName == "tags" {
							newsArticle.Keywords = tagVal.([]string)
							for i := range newsArticle.Keywords {
								newsArticle.Keywords[i] = strings.TrimSpace(newsArticle.Keywords[i])
							}
						}

						if tagName == "section" {
							newsArticle.ArticleSection = tagVal.(string)
						}
					}
				}
			}

			items := aggr.crawl.MicroData.(map[string]interface{})["items"].([]interface{})
			if len(items) >= 1 {
				for _, item := range items {
					itemMap := item.(map[string]interface{})
					for k, v := range itemMap {
						if k == "properties" {
							properties := v.(map[string]interface{})

							if dateList, ok := properties["datePublished"]; ok {
								date := dateList.([]interface{})[0]
								t, err := dateparse.ParseAny(date.(string))
								if err == nil {
									newsArticle.DatePublished = t
								} else {
									log.Println(err)
								}
							}

							if dateList, ok := properties["dateModified"]; ok {
								date := dateList.([]interface{})[0]
								t, err := dateparse.ParseAny(date.(string))
								if err == nil {
									newsArticle.DateModified = t
								}
							}

							if headline, ok := properties["alternativeHeadline"]; ok {
								h := headline.([]interface{})[0]
								newsArticle.AlternativeHeadline = h.(string)
							}

							if headline, ok := properties["articleSection"]; ok {
								h := headline.([]interface{})[0]
								newsArticle.ArticleSection = h.(string)
							}
						}
					}
				}
			}

			for _, data := range aggr.crawl.JSONData {
				if mdata, ok := data.(map[string]interface{}); ok {
					if articleSection, ok := mdata["articleSection"]; ok {
						if as, ok := articleSection.(string); ok && as != "" {
							newsArticle.ArticleSection = as
						}
					}

					if alternativeHeadline, ok := mdata["alternativeHeadline"]; ok {
						if as, ok := alternativeHeadline.(string); ok && as != "" {
							newsArticle.AlternativeHeadline = as
						}
					}

					if datei, ok := mdata["datePublished"]; ok {
						date := datei.(string)
						t, err := dateparse.ParseAny(date)
						if err == nil {
							newsArticle.DatePublished = t
						}
					}

					if datei, ok := mdata["dateModified"]; ok {
						date := datei.(string)
						t, err := dateparse.ParseAny(date)
						if err == nil {
							newsArticle.DateModified = t
						}
					}

					if thumbi, ok := mdata["thumbnailUrl"]; ok {
						newsArticle.Images = append([]string{thumbi.(string)}, newsArticle.Images...)
					}
				}
			}

			newsArticle.WordCount = len(strings.Fields(newsArticle.ArticleBody))
			entities := a.processEntities(newsArticle)

			aggr.article = newsArticle
			aggr.entities = entities
		}
	}
}

func (a *Aggregator) processEntities(na NewsArticle) []tag {
	body := na.ArticleBody
	title := na.Headline
	desc := na.Description
	tagMap := map[string]tag{}

	tags := a.extractEntities(body)
	for k, u := range tags {
		u.FromBody = true
		tagMap[k] = u
	}

	tags = a.extractEntities(title)
	for k, u := range tags {
		if _, ok := tagMap[k]; !ok {
			tagMap[k] = u
		}
		g := tagMap[k]
		g.FromTitle = true
		tagMap[k] = g
	}

	tags = a.extractEntities(desc)
	for k, u := range tags {
		if _, ok := tagMap[k]; !ok {
			tagMap[k] = u
		}
		g := tagMap[k]
		g.FromDescription = true
		tagMap[k] = g
	}

	tagList := []tag{}

	for _, u := range tagMap {
		tagList = append(tagList, u)
	}

	return tagList
}

func (a *Aggregator) extractEntities(txt string) map[string]tag {
	tokens := ner.Tokenize(txt)
	es, err := a.ne.Extract(tokens)
	if err != nil {
		return nil
	}

	tags := make(map[string]tag)

	for _, v := range es {
		reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
		if err != nil {
			log.Fatal(err)
		}
		name := v.Name
		name = reg.ReplaceAllString(name, "")

		if name == "" {
			continue
		}

		_, exists := tags[name]
		if !exists {
			tags[name] = tag{
				Name: name,
				Kind: a.ne.Tags()[v.Tag],
			}
		}

		g := tags[name]
		g.Count++
		g.Total += v.Score
		tags[name] = g
	}

	for k, v := range tags {
		if v.Total < 0.67 {
			delete(tags, k)
		}
	}

	return tags
}

func (a *Aggregator) saveAggregate(aggr *Aggregation) error {
	for _, tag := range aggr.keywords {
		err := a.saveKeyWord(tag, aggr.documentRef)
		if err != nil {
			err = errors.Wrap(err, "Unable to save keyword {"+tag+"}")
			log.Println(err)
		}
	}

	for _, card := range aggr.cards {
		err := a.saveCard(card, aggr.documentRef)
		if err != nil {
			err = errors.Wrap(err, "Unable to save keyword {"+card["type"].(string)+"}")
			log.Println(err)
		}
	}

	naRef, err := a.saveArticle(aggr.article)
	if err != nil {
		return errors.Wrap(err, "unable to save article")
	}

	// Save each entity
	for _, e := range aggr.entities {
		e.PageRef = aggr.documentRef
		err := a.saveEntity(e, naRef)
		if err != nil {
			err = errors.Wrap(err, "unabelt to save entity {"+e.Name+"}")
			log.Println(err)
			return err
		}
	}

	return nil
}

func (a *Aggregator) saveKeyWord(key string, ref string) error {
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		return err
	}

	key = strings.TrimSpace(key)
	key = strings.ToLower(key)
	keyref := reg.ReplaceAllString(key, "")
	keyref = strcase.ToSnake(keyref)

	kw := keyword{
		Keyword: key,
		Key:     keyref,
	}

	exists, err := a.namedKeywords.DocumentExists(context.Background(), keyref)
	if err != nil {
		log.Println(kw.Keyword, err)
	}

	var m driver.DocumentMeta
	if !exists {
		m, err = a.namedKeywords.CreateDocument(context.Background(), kw)
		if err != nil {
			log.Println(kw.Keyword, err)
		}
	} else {
		var k keyword
		m, err = a.namedKeywords.ReadDocument(context.Background(), keyref, &k)
		if err != nil {
			log.Println(kw.Keyword, err)
		}
	}

	_, err = a.keyRef.CreateDocument(context.Background(), map[string]string{
		"_to":   m.ID.String(),
		"_from": ref,
		"name":  "has_keyword",
		"date":  time.Now().Format(time.RFC3339),
	})

	if err != nil {
		log.Println(kw.Keyword, err)
		return err
	}

	return nil
}

func (a *Aggregator) saveArticle(na NewsArticle) (string, error) {
	m, ok := a.articleExists(na.Headline)
	if !ok {
		m2, err := a.content.CreateDocument(context.Background(), na)
		if !driver.IsArangoErrorWithErrorNum(err, 1210) && err != nil {
			log.Println(na, err)
			return "", err
		}
		m = &m2
	}

	ref := m.ID.String()

	_, err := a.pageToContent.CreateDocument(context.Background(), map[string]string{
		"_from": na.PageRef,
		"_to":   ref,
		"type":  "references",
	})

	if !driver.IsArangoErrorWithErrorNum(err, 1210) && err != nil {
		log.Println(na.PageRef, "content/"+ref, na.Headline, err)
		return "", err
	}

	return ref, nil
}

func (a *Aggregator) saveEntity(entity tag, content_ref string) error {
	m, exists := a.entityExists(entity.Name)
	if !exists {
		m2, err := a.entities.CreateDocument(context.Background(), entity)
		if !driver.IsArangoErrorWithErrorNum(err, 1210) && err != nil {
			log.Println(entity.Name, err)
			return err
		}
		m = &m2
	}

	_, err := a.pageToEntity.CreateDocument(context.Background(), map[string]interface{}{
		"_from":           entity.PageRef,
		"_to":             m.ID.String(),
		"count":           entity.Count,
		"total":           entity.Total,
		"fromBody":        entity.FromBody,
		"fromTitle":       entity.FromTitle,
		"fromDescription": entity.FromDescription,
		"type":            "references",
	})

	if !driver.IsArangoErrorWithErrorNum(err, 1210) && err != nil {
		log.Println(entity.Name, err)
		return err
	}

	_, err = a.entityRef.CreateDocument(context.Background(), map[string]interface{}{
		"_from":           m.ID.String(),
		"_to":             content_ref,
		"count":           entity.Count,
		"total":           entity.Total,
		"fromBody":        entity.FromBody,
		"fromTitle":       entity.FromTitle,
		"fromDescription": entity.FromDescription,
		"type":            "references",
	})

	if !driver.IsArangoErrorWithErrorNum(err, 1210) && err != nil {
		log.Println(entity.Name, err)
		return err
	}

	return nil
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

type document struct {
	URL         string `json:"URL"`
	Title       string
	Description string
	Data        interface{}
	JData       interface{}
	MData       interface{}
}

type NewsArticle struct {
	ArticleBody         string
	ArticleSection      string
	WordCount           int
	DatePublished       time.Time
	DateModified        time.Time
	Headline            string
	AlternativeHeadline string
	Description         string
	Keywords            []string
	PageRef             string
	Images              []string
}

type tag struct {
	Name            string
	Kind            string
	Count           int     `json:"-"`
	Total           float64 `json:"-"`
	FromBody        bool    `json:"-"`
	FromTitle       bool    `json:"-"`
	FromDescription bool    `json:"-"`
	PageRef         string  `json:"-"`
}

func (a *Aggregator) saveCard(card map[string]interface{}, ref string) error {
	m, err := a.meta.CreateDocument(context.Background(), card)
	if err != nil {
		log.Println(ref, card, err)
		return err
	}

	_, err = a.metaRef.CreateDocument(context.Background(), map[string]string{
		"_from": ref,
		"_to":   m.ID.String(),
		"type":  card["type"].(string),
	})

	return err
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

func (a *Aggregator) articleExists(k string) (*driver.DocumentMeta, bool) {
	query := "FOR d IN @col FILTER d.Headline == @headline LIMIT 1 RETURN d"
	cursor, err := a.db.Query(context.Background(), query, map[string]interface{}{
		"headline": k,
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

func (a *Aggregator) entityExists(k string) (*driver.DocumentMeta, bool) {
	query := "FOR d IN entities FILTER d.Name == @name LIMIT 1 RETURN d"
	cursor, err := a.db.Query(context.Background(), query, map[string]interface{}{
		"name": k,
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
