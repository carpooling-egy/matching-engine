package Publisher

import "matching-engine/internal/model"

type MatchingResultSink interface {
	// Add adds a new matching result to the sink
	publish(matchingResults []model.MatchingResult) error
}
