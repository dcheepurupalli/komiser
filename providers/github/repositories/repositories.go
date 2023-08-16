package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

		variables, _, err := client.GithubClient.Actions.ListRepoVariables(ctx, client.Name, *repository.Name, nil)
		if err != nil {
			return resources, err
		}

		sbom, err := GetRepoDependencyGraphSBOM(ctx, client, repository)
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

		// Convert variables to key-value pairs
		variablePairs := make([]KeyValuePair, 0)
		for _, variable := range variables.Variables {
			variablePairs = append(variablePairs, KeyValuePair{
				Key:   variable.Name,
				Value: variable.Value,
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
			Variables:  variablePairs,
			SBOM:       sbom,
		})
	}

	return resources, nil
}

func sendRequest(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set request headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetRepoDependencyGraphSBOM(ctx context.Context, client providers.ProviderClient, repository *github.Repository) (string, error) {
	// Set the necessary headers
	headers := map[string]string{
		"Accept":        "application/vnd.github+json",
		"Authorization": "Bearer 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",
		"X-GitHub-Api":  "2022-11-28",
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/dependency-graph/sbom", client.Name, *repository.Name)

	// Send GET request
	resp, err := sendRequest("GET", url, headers, nil)
	if err != nil {
		log.Fatalln("Request error:", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Println("API request failed with status code:", resp.StatusCode)
		return "", err
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}
	return string(body), nil
}
