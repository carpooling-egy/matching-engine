package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/service/checker"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterCheckers registers checking services
func RegisterCheckers(c *dig.Container) {
	utils.Must(c.Provide(checker.NewCapacityChecker, dig.Name("capacity_checker")))
	utils.Must(c.Provide(checker.NewOverlapChecker, dig.Name("overlap_checker")))
	utils.Must(c.Provide(checker.NewDetourTimeChecker, dig.Name("detour_checker")))
	utils.Must(c.Provide(checker.NewPreferenceChecker, dig.Name("preference_checker")))
	utils.Must(c.Provide(provideCompositeChecker))
}

// CheckerParams contains all checkers for the composite checker
type CheckerParams struct {
	dig.In

	CapacityChecker   checker.Checker `name:"capacity_checker"`
	DetourTimeChecker checker.Checker `name:"detour_checker"`
	OverlapChecker    checker.Checker `name:"overlap_checker"`
	PreferenceChecker checker.Checker `name:"preference_checker"`
}

// provideCompositeChecker provides a composite checker with all other checkers
func provideCompositeChecker(params CheckerParams) checker.Checker {
	return checker.NewCompositeChecker(
		params.OverlapChecker,
		params.CapacityChecker,
		params.PreferenceChecker,
		params.DetourTimeChecker,
	)
}
