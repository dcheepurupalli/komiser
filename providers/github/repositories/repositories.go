package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
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

		secrets, _, err := client.GithubClient.Actions.ListRepoSecrets(ctx, client.Name, *repository.Name, nil)
		if err != nil {
			return resources, err
		}

		// Convert secrets to key-value pairs
		secretPairs := make([]KeyValuePair, 0)
		for _, secret := range secrets.Secrets {
			secretPairs = append(secretPairs, KeyValuePair{
				Key: secret.Name,
			})
		}
		tags := make([]Tag, 0)
		for _, tag := range repository.Topics {
			fmt.Println("Tags", tag)
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
			Secrets:    secretPairs,
		})
	}

	return resources, nil
}
