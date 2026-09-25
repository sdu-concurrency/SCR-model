package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("questions")
		if err != nil {
			return err
		}

		// Add the optional job_function_schema JSON field if missing.
		if collection.Fields.GetByName("job_function_schema") == nil {
			collection.Fields.Add(&core.JSONField{
				Name: "job_function_schema",
			})
			if err := app.Save(collection); err != nil {
				return err
			}
		}

		// Backfill: seed the job function options for existing question sets
		// without one, so the survey renders them instead of the hardcoded
		// fallback options.
		records, err := app.FindAllRecords("questions")
		if err != nil {
			return err
		}
		for _, record := range records {
			current := record.GetString("job_function_schema")
			if current == "" || current == "null" {
				record.Set("job_function_schema", jobFunctionSchemaDefault)
				if err := app.Save(record); err != nil {
					return err
				}
			}
		}

		return nil
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("questions")
		if err != nil {
			return err
		}

		// Down: clear the seeded values, then remove the field.
		records, err := app.FindAllRecords("questions")
		if err != nil {
			return err
		}
		for _, record := range records {
			if record.Get("job_function_schema") != nil {
				record.Set("job_function_schema", nil)
				if err := app.Save(record); err != nil {
					return err
				}
			}
		}

		collection.Fields.RemoveByName("job_function_schema")
		return app.Save(collection)
	})
}

// jobFunctionSchemaDefault contains the multilingual job function options for
// the logistics question set. Same structure as the other schemas: a single
// category whose items carry a "value" and a multilingual "label".
var jobFunctionSchemaDefault = []map[string]any{
	{
		"label": map[string]string{"en": "Job function", "da": "Jobfunktion"},
		"items": []map[string]any{
			{"value": "Salg og Customer Service", "label": map[string]string{"en": "Sales and Customer Service", "da": "Salg og Customer Service"}},
			{"value": "Spedition og Transport Operations", "label": map[string]string{"en": "Freight forwarding and transport operations", "da": "Spedition og Transport Operations"}},
			{"value": "Transportindkøb og Carrier Management", "label": map[string]string{"en": "Transport procurement and carrier management", "da": "Transportindkøb og Carrier Management"}},
			{"value": "Terminal, lager og fysisk godshåndtering", "label": map[string]string{"en": "Terminal, warehouse and physical goods handling", "da": "Terminal, lager og fysisk godshåndtering"}},
			{"value": "Told, dokumentation og compliance", "label": map[string]string{"en": "Customs, documentation and compliance", "da": "Told, dokumentation og compliance"}},
			{"value": "IT og digitale systemer", "label": map[string]string{"en": "IT and digital systems", "da": "IT og digitale systemer"}},
			{"value": "Finans og administration", "label": map[string]string{"en": "Finance and administration", "da": "Finans og administration"}},
			{"value": "Andet", "label": map[string]string{"en": "Other", "da": "Andet"}},
		},
	},
}
