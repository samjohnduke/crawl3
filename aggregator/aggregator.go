package reporter

import "context"

// Service implementations provide a view over the entire dataset of the application
type Service interface {
	Query(queryStr string, ctx *context.Context) (result interface{}, err error)
	QueryWithUpdater()
}

// Client implementations provide a way for others to read the data
type Client interface {
	Query(queryStr string, ctx *context.Context) (result interface{}, err error)
}

// Aggregator is responsible for outputing the data requested by a client
type Aggregator struct {
}
