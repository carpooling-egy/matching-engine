package publisher

import (
	"matching-engine/internal/model"
)

// Publisher represents a messaging system that can publish messages
type Publisher interface {
	// Publish the matching results to the messaging system
	Publish(results []*model.MatchingResult) error

	// Close releases resources used by the publisher
	Close() error
}
