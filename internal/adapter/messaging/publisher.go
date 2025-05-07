package messaging

import (
	"matching-engine/internal/model"
)

// Publisher represents a messaging system that can publish messages
type Publisher interface {
	// PublishMatchingResults publish the matching results to the messaging system
	PublishMatchingResults(results []*model.MatchingResult) error

	// Close releases resources used by the publisher
	Close() error
}
