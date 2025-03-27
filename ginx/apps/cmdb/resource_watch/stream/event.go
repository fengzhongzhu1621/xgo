package stream

import "go.mongodb.org/mongo-driver/mongo"

// Event TODO
type Event struct {
	database string
	client   *mongo.Client
}

// NewEvent TODO
func NewEvent(client *mongo.Client, db string) (*Event, error) {
	return &Event{client: client, database: db}, nil
}
