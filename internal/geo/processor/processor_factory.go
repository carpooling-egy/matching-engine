package processor

import (
	"matching-engine/internal/model"
)

// ProcessorFactory defines the interface for creating GeospatialProcessor instances.
type ProcessorFactory interface {
	CreateProcessor(offer *model.Offer) (GeospatialProcessor, error)
}
