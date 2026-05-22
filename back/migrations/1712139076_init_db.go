package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	// Initialise the database schema.
	// No superuser is created here — use the PocketBase setup wizard at /_/
	// on first startup to create your admin account.
	m.Register(
		func(app core.App) error {
			return nil
		}, func(app core.App) error {
			return nil
		})
}
