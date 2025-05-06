package publisher

import "matching-engine/internal/model"

type ResultSink interface {
	// Add adds a new matching result to the sink
	publish(matchingResults []model.MatchingResult) error
}
