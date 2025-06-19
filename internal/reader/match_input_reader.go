package reader

import (
	"context"
	"matching-engine/internal/model"
)

// MatchInputReader defines the interface for reading offers and requests from an input source
type MatchInputReader interface {
	GetOffersAndRequests(ctx context.Context) ([]*model.Request, []*model.Offer, bool, error)
	Close() error
}
