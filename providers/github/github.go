package github

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/tailwarden/komiser/providers"
	"github.com/tailwarden/komiser/providers/github/repositories"
	"github.com/tailwarden/komiser/utils"
	"github.com/uptrace/bun"
)

func listOfSupportedServices() []providers.FetchDataFunction {
	return []providers.FetchDataFunction{
		repositories.Repositories,
	}
}

func FetchResources(ctx context.Context, client providers.ProviderClient, db *bun.DB, telemetry bool, analytics utils.Analytics) {
	for _, fetchResources := range listOfSupportedServices() {
		resources, err := fetchResources(ctx, client)
		if err != nil {
			log.Printf("[%s][Github] %s", client.Name, err)
		} else {
			for _, resource := range resources {
				_, err := db.NewInsert().Model(&resource).On("CONFLICT (resource_id) DO UPDATE").Set("cost = EXCLUDED.cost").Exec(context.Background())
				if err != nil {
					logrus.WithError(err).Errorf("db trigger failed")
				}
			}
			if telemetry {
				analytics.TrackEvent("discovered_resources", map[string]interface{}{
					"provider":  "Github",
					"resources": len(resources),
				})
			}
		}
	}
}
