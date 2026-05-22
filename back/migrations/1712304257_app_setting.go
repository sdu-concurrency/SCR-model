package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		settings := app.Settings()
		settings.Meta.AppName = "Supply Chain Resilience - Process model"
		settings.Logs.MaxDays = 2

		return app.Save(settings)
	}, nil)
}
