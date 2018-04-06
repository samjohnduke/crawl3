package schedular

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// SQLGormStorage implements the schedular.Store interface using the GORM library connected
// a SQL database.
type SQLGormStorage struct {
	db *gorm.DB
}

// NewSQLGormStore builds a new store from a database, first creating the the tables
// -- Visit
// -- Queued
// -- QueueItem
// if they aren't at the latest version
func NewSQLGormStore(db *gorm.DB) (Store, error) {
	db.AutoMigrate(&Visit{})
	db.AutoMigrate(&Queued{})
	db.AutoMigrate(&QueueItem{})

	return &SQLGormStorage{db: db}, nil
}

// Queue pushes the provided url into the database
func (s *SQLGormStorage) Queue(u string) error {
	if err := s.db.Create(&QueueItem{
		URL: u,
	}).Error; err != nil {
		return errors.Wrapf(err, "unable to push {%s} into the queue", u)
	}

	return nil
}

// QueueAt pushes the provided url into the queue at some future time
func (s *SQLGormStorage) QueueAt(u string, t time.Time) error {
	if err := s.db.Create(&Queued{
		URL: u,
		At:  t,
	}).Error; err != nil {
		return errors.Wrapf(err, "unable to push {%s} at {%t}", u, t)
	}
	return nil
}

// IsQueued checks if the url is in the queue
func (s *SQLGormStorage) IsQueued(u string) bool {
	var queue QueueItem
	if s.db.Where(&QueueItem{URL: u}).First(&queue).RecordNotFound() {
		return false
	}
	return true
}

// Visit updates the databse to ensure that the url is is in the visit table
// with the latest hash and an updated count
func (s *SQLGormStorage) Visit(u string, hash string) (Visit, error) {
	var visit Visit
	now := time.Now()

	if s.db.Where(&Visit{URL: u}).First(&visit).RecordNotFound() {
		visit = Visit{
			URL:             u,
			UpdateFrequency: 15 * time.Minute,
			LastUpdate:      &now,
			VisitCount:      1,
			LastHash:        hash,
		}

		if err := s.db.Create(&visit).Error; err != nil {
			return Visit{}, errors.Wrap(err, "Unable to create visit")
		}
	} else {
		visit.LastUpdate = &now
		visit.VisitCount++
		visit.LastHash = hash

		if err := s.db.Save(&visit).Error; err != nil {
			return Visit{}, errors.Wrap(err, "Unable to update visit")
		}
	}

	if err := s.db.Delete(&QueueItem{URL: u}).Error; err != nil {
		return Visit{}, errors.Wrap(err, "Unable to remove from the queue")
	}

	return visit, nil
}

// ShouldVisit determines if it is appropriate to push the URL into the queue
func (s *SQLGormStorage) ShouldVisit(u string) bool {
	return !s.HasVisited(u) && !s.IsQueued(u)
}

// HasVisited checks if the URL is in the database as a previous visit
func (s *SQLGormStorage) HasVisited(u string) bool {
	var visit Visit
	if s.db.Where(&Visit{URL: u}).First(&visit).RecordNotFound() {
		return false
	}
	return true
}

// Reschedule will allow the store to push items in the queue back into the schedular
func (s *SQLGormStorage) Reschedule(service Service) {

}
