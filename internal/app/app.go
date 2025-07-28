package app

import (
	"context"
	"go.uber.org/dig"
	"matching-engine/internal/app/di"
	"matching-engine/internal/app/starter"
)

// App represents the application and its dependencies
type App struct {
	container *dig.Container
}

// NewApp creates and initializes a new application
func NewApp() *App {
	return &App{
		container: di.BuildContainer(),
	}
}

// Run starts the application
func (app *App) Run(ctx context.Context) error {
	return app.container.Invoke(func(s *starter.StarterService) error {
		return s.Start(ctx)
	})
}
