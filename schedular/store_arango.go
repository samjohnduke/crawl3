package schedular

import (
	"context"
	"time"

	driver "github.com/arangodb/go-driver"
	"github.com/pkg/errors"
)

type ArangoStorage struct {
	db         driver.Database
	queue      driver.Collection
	queueAfter driver.Collection
	visits     driver.Collection
}

func NewArangoStore(db driver.Database) (Store, error) {

	queue, err := db.Collection(context.Background(), "queue")
	if err != nil {
		return nil, err
	}

	queueAfter, err := db.Collection(context.Background(), "queue_after")
	if err != nil {
		return nil, err
	}

	visited, err := db.Collection(context.Background(), "visits")
	if err != nil {
		return nil, err
	}

	return &ArangoStorage{
		db:         db,
		queue:      queue,
		queueAfter: queueAfter,
		visits:     visited,
	}, nil
}

func (s *ArangoStorage) Queue(u string) error {
	v := Visit{
		URL: u,
	}

	_, err := s.queue.CreateDocument(context.Background(), v)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArangoStorage) QueueAt(u string, t time.Time) error {
	v := Queued{URL: u, At: t}

	_, err := s.queueAfter.CreateDocument(context.Background(), v)
	if err != nil {
		return err
	}

	return nil
}

func (s *ArangoStorage) IsQueued(u string) bool {
	query := "FOR d IN queue FILTER d.URL == @url LIMIT 1 RETURN d"
	cursor, err := s.db.Query(context.Background(), query, map[string]interface{}{
		"url": u,
	})
	if err != nil {
		return false
	}
	defer cursor.Close()

	var doc Visit
	_, err = cursor.ReadDocument(context.Background(), &doc)
	if driver.IsNoMoreDocuments(err) {
		return false
	} else if err != nil {
		return false
	}

	return true
}

func (s *ArangoStorage) Visit(u string, hash string) (Visit, error) {
	query := "FOR d IN visits FILTER d.URL == @url LIMIT 1 RETURN d"
	cursor, err := s.db.Query(context.Background(), query, map[string]interface{}{
		"url": u,
	})
	if err != nil {
		return Visit{}, errors.WithStack(err)
	}
	defer cursor.Close()

	var doc Visit
	m, err := cursor.ReadDocument(context.Background(), &doc)
	if driver.IsNoMoreDocuments(err) {
		doc = Visit{
			URL:             u,
			UpdateFrequency: 15 * time.Minute,
		}
	} else if err != nil {
		return Visit{}, errors.Wrap(err, "unable to read last visit")
	}

	now := time.Now()
	doc.LastUpdate = &now
	doc.VisitCount++
	doc.LastHash = hash

	if err == nil {
		_, err = s.visits.ReplaceDocument(context.Background(), m.Key, doc)
		if err != nil {
			return Visit{}, errors.Wrap(err, "unabel to replace document")
		}
	} else {
		_, err = s.visits.CreateDocument(context.Background(), doc)
		if err != nil {
			return Visit{}, errors.Wrap(err, "unable to create document")
		}
	}

	query2 := "FOR d IN queue FILTER d.URL == @url LIMIT 1 RETURN d"
	cursor2, err := s.db.Query(context.Background(), query2, map[string]interface{}{
		"url": u,
	})
	if err != nil {
		return Visit{}, errors.Wrap(err, "unable to query the queue")
	}
	defer cursor2.Close()

	var doc2 Visit
	m, err = cursor2.ReadDocument(context.Background(), &doc2)
	if driver.IsNoMoreDocuments(err) {
		return Visit{}, errors.Wrap(err, "url not in queue")
	} else if err != nil {
		return Visit{}, errors.Wrap(err, "unable to get the url from the queue")
	}

	s.queue.RemoveDocument(context.Background(), m.Key)

	return Visit{}, nil
}

func (s *ArangoStorage) ShouldVisit(u string) bool {
	var visit Visit
	var shouldVisit = true

	query := "FOR d IN visits FILTER d.URL == @url LIMIT 1 RETURN d"
	cursor, err := s.db.Query(context.Background(), query, map[string]interface{}{
		"url": u,
	})
	if err != nil {
		return false
	}
	defer cursor.Close()

	_, err = cursor.ReadDocument(context.Background(), &visit)
	if driver.IsNoMoreDocuments(err) {
		shouldVisit = true
	} else if err != nil {
		return false
	}

	if visit.NextUpdate.After(time.Now()) {
		shouldVisit = false
	}

	return shouldVisit
}

func (s *ArangoStorage) HasVisited(u string) bool {
	var hasVisited = true

	query := "FOR d IN visits FILTER d.URL == @url LIMIT 1 RETURN d"
	cursor, err := s.db.Query(context.Background(), query, map[string]interface{}{
		"url": u,
	})
	if err != nil {
		return false
	}
	defer cursor.Close()

	var visit Visit
	_, err = cursor.ReadDocument(context.Background(), &visit)
	if driver.IsNoMoreDocuments(err) {
		hasVisited = false
	} else if err != nil {
		return true
	}

	return hasVisited
}

func (s *ArangoStorage) Reschedule(Service) {

}
