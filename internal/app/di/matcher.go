package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/service/checker"
	"matching-engine/internal/service/earlypruning"
	"matching-engine/internal/service/matcher"
	"matching-engine/internal/service/matchevaluator"
	"matching-engine/internal/service/maximummatching"
	"matching-engine/internal/service/pathgeneration/planner"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterMatchingServices registers matching algorithm services
func RegisterMatchingServices(c *dig.Container) {
	utils.Must(c.Provide(provideMatchEvaluator))
	utils.Must(c.Provide(earlypruning.NewPreChecksCandidateGenerator))
	utils.Must(c.Provide(maximummatching.NewHopcroftKarp))
	utils.Must(c.Provide(matcher.NewMatcher))
}

// MatchEvaluatorParams contains the dependencies for the match evaluator
type MatchEvaluatorParams struct {
	dig.In

	PathPlanner       planner.PathPlanner
	PreferenceChecker checker.Checker `name:"preference_checker"`
}

// provideMatchEvaluator provides a match evaluator
func provideMatchEvaluator(params MatchEvaluatorParams) matchevaluator.Evaluator {
	return matchevaluator.NewMatchEvaluator(params.PathPlanner, params.PreferenceChecker)
}
