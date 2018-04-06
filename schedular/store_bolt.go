package schedular

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/boltdb/bolt"
)

type BoltStorage struct {
	db *bolt.DB
}

func NewStore(db *bolt.DB) (Store, error) {
	err := db.Update(func(tx *bolt.Tx) error {

		_, err := tx.CreateBucketIfNotExists([]byte("Visits"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("Queue"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("QueueAfter"))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &BoltStorage{
		db: db,
	}, nil
}

func (s *BoltStorage) Queue(u string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		q := tx.Bucket([]byte("Queue"))

		visit := Visit{URL: u}
		visitJ, err := json.Marshal(visit)
		if err != nil {
			return err
		}

		q.Put([]byte(u), visitJ)

		return nil
	})

	return err
}

func (s *BoltStorage) QueueAt(u string, t time.Time) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		q := tx.Bucket([]byte("QueueAfter"))

		visit := Queued{URL: u, At: t}
		visitJ, err := json.Marshal(visit)
		if err != nil {
			return err
		}

		q.Put([]byte(u), visitJ)

		return nil
	})
	return err
}

func (s *BoltStorage) Visit(u string, hash string) (Visit, error) {
	var visit Visit
	err := s.db.Update(func(tx *bolt.Tx) error {
		q := tx.Bucket([]byte("Visits"))

		var v Visit
		now := time.Now()
		kv := q.Get([]byte(u))
		if kv == nil {
			v = Visit{
				URL:             u,
				LastUpdate:      &now,
				UpdateFrequency: 15 * time.Minute,
				VisitCount:      1,
				LastHash:        hash,
			}
		} else {
			err := json.Unmarshal(kv, &v)
			if err != nil {
				log.Println(err)
				return err
			}

			v.LastUpdate = &now
			v.VisitCount++

			if v.LastHash == hash && v.UpdateBackoff <= 12 {
				v.UpdateFrequency = time.Duration(int(math.Pow(2, float64(v.UpdateBackoff)))) * 15 * time.Minute
				v.UpdateBackoff++
				now := time.Now().Add(v.UpdateFrequency)
				v.NextUpdate = &now
			} else if v.LastHash != hash {
				v.UpdateFrequency = 15 * time.Minute
				v.UpdateBackoff = 1
			}
		}

		visit = v

		visitJ, err := json.Marshal(v)
		if err != nil {
			log.Println(err)
			return err
		}

		q.Put([]byte(u), visitJ)

		q = tx.Bucket([]byte("Queue"))

		err = q.Delete([]byte(u))
		if err != nil {
			log.Println(err)
			return err
		}

		return nil
	})

	return visit, err
}

func (s *BoltStorage) ShouldVisit(u string) bool {
	var visit Visit
	var shouldVisit = true

	err := s.db.Update(func(tx *bolt.Tx) error {
		q := tx.Bucket([]byte("Visits"))

		var v Visit
		kv := q.Get([]byte(u))
		if kv != nil {
			err := json.Unmarshal(kv, &v)
			if err != nil {
				log.Println(err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		return false
	}

	if visit.NextUpdate.After(time.Now()) {
		shouldVisit = false
	}

	return shouldVisit
}

func (s *BoltStorage) HasVisited(u string) bool {

	success := false

	s.db.Update(func(tx *bolt.Tx) error {
		q := tx.Bucket([]byte("Visits"))

		v := q.Get([]byte(u))
		if v == nil {
			return nil
		}

		success = true
		return nil
	})

	return success
}

func (s *BoltStorage) IsQueued(u string) bool {
	queued := false

	s.db.View(func(tx *bolt.Tx) error {
		q := tx.Bucket([]byte("Queue"))

		v := q.Get([]byte(u))
		if v == nil {
			return nil
		}

		queued = true
		return nil
	})

	return queued
}

func (s *BoltStorage) Reschedule(sched Service) {
	err := s.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("QueueAfter"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}

			if t.Before(time.Now()) {
				var vv Queued
				err := json.Unmarshal(v, vv)
				if err != nil {
					return err
				}

				s.Queue(vv.URL)
				sched.Schedule(vv.URL)
				c.Delete()
			} else {
				break
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
