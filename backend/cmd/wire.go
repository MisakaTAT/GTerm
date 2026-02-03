//go:build wireinject

package cmd

import (
	"github.com/MisakaTAT/GTerm/backend/initialize"
	"github.com/MisakaTAT/GTerm/backend/services"
	"github.com/google/wire"
)

func NewApp() *App {
	wire.Build(
		AppSet,
		initialize.InitDatabase,
		initialize.InitHTTPServer,
		services.ProviderSet,
	)
	return new(App)
}
