package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/tailwarden/komiser/models"
	. "github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"
)

func Repositories(ctx context.Context, client providers.ProviderClient) ([]models.Resource, error) {
	resources := make([]models.Resource, 0)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{},
	}

	repositories, _, err := client.GithubClient.Repositories.ListByOrg(ctx, client.Name, opt)
	if err != nil {
		return resources, err
	}

	for _, repository := range repositories {

		tags := make([]Tag, 0)
		for _, tag := range repository.Topics {
			if strings.Contains(tag, ":") {
				parts := strings.Split(tag, ":")
				tags = append(tags, models.Tag{
					Key:   parts[0],
					Value: parts[1],
				})
			} else {
				tags = append(tags, models.Tag{
					Key:   tag,
					Value: tag,
				})
			}
		}

		jsonData, err := json.Marshal(repository)
		if err != nil {
			fmt.Println("Error:", err)
		}
		jsonString := string(jsonData)

		resources = append(resources, models.Resource{
			Provider:   "Github",
			Account:    client.Name,
			Service:    "Repository",
			ResourceId: *repository.Name,
			Cost:       0,
			Name:       *repository.Name,
			FetchedAt:  time.Now(),
			CreatedAt:  repository.CreatedAt.Time,
			Tags:       tags,
			Data:       jsonString,
			Link:       *repository.URL,
		})
	}

	return resources, nil
}
