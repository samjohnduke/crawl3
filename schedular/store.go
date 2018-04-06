package schedular

import (
	"time"
)

type Store interface {
	Queue(string) error
	QueueAt(string, time.Time) error
	IsQueued(string) bool
	Visit(string, string) (Visit, error)
	ShouldVisit(string) bool
	HasVisited(string) bool
	Reschedule(Service)
}

type Visit struct {
	URL             string `gorm:"primary_key"`
	LastVisit       *time.Time
	LastUpdate      *time.Time
	LastHash        string
	UpdateFrequency time.Duration
	UpdateBackoff   int64
	NextUpdate      *time.Time
	VisitCount      int64
}

type Queued struct {
	URL string `gorm:"primary_key"`
	At  time.Time
}

type QueueItem struct {
	URL string `gorm:"primary_key"`
}
