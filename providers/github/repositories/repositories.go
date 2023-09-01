package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v54/github"
	"github.com/tailwarden/komiser/models"
	. "github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"
)

type Gitleaks struct {
	Description string   `json:"Description"`
	StartLine   int      `json:"StartLine"`
	EndLine     int      `json:"EndLine"`
	StartColumn int      `json:"StartColumn"`
	EndColumn   int      `json:"EndColumn"`
	Match       string   `json:"Match"`
	Secret      string   `json:"Secret"`
	File        string   `json:"File"`
	SymlinkFile string   `json:"SymlinkFile"`
	Commit      string   `json:"Commit"`
	Entropy     float64  `json:"Entropy"`
	Author      string   `json:"Author"`
	Email       string   `json:"Email"`
	Date        string   `json:"Date"`
	Message     string   `json:"Message"`
	Tags        []string `json:"Tags"`
	RuleID      string   `json:"RuleID"`
	Fingerprint string   `json:"Fingerprint"`
}

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

		// SBOM Changes
		sbom, _, err := client.GithubClient.DependencyGraph.GetSBOM(ctx, client.Name, *repository.Name)
		if err != nil {
			fmt.Println("Error getting SBOM:", err)
			// return resources, err
		}
		sbomJson, err := json.MarshalIndent(sbom, "", "    ")
		if err != nil {
			fmt.Println("Error:", err)
		}
		sbomString := string(sbomJson)
		// fmt.Println(sbomString)

		// Git Leaks Cloning repository
		err = os.MkdirAll("./repos", os.ModePerm)
		if err != nil {
			fmt.Println("Error creating repos directory:", err)
			return resources, err
		}

		gitLeaks, err := cloneRepositoryInTemp(ctx, client, repository)
		if err != nil {
			return resources, err
		}
		gitLeaksJson, err := json.Marshal(gitLeaks)
		if err != nil {
			fmt.Println("Error:", err)
		}
		gitLeaksString := string(gitLeaksJson)

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
			SBOM:       sbomString,
			Unmanaged:  gitLeaksString,
		})
	}

	return resources, nil
}

// Code is written temporary not following best practices
func cloneRepositoryInTemp(ctx context.Context, client providers.ProviderClient, repository *github.Repository) ([]Gitleaks, error) {
	repoName := *repository.Name
	repoURL := *repository.CloneURL

	// Cloning the repository
	fmt.Printf("Cloning repository %s...\n", repoName)
	cmd := exec.Command("git", "clone", repoURL)
	cmd.Dir = "./repos"
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error cloning repository:", err)
	}

	// Running gitleaks
	fmt.Printf("Running gitleaks on repository %s...\n", repoName)
	cmd = exec.Command("gitleaks", "detect", "-s", ".", "-r", "gitleaks-report.json", "-f", "json")
	cmd.Dir = fmt.Sprintf("./repos/%s", repoName)
	// resolve below error
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running gitleaks:", err)
	}

	// Read the captured JSON report
	reportBytes, err := ioutil.ReadFile(fmt.Sprintf("./repos/%s/gitleaks-report.json", repoName))
	if err != nil {
		fmt.Println("Error reading gitleaks report:", err)
	}

	// Parse the JSON report
	var gitleaks []Gitleaks
	err = json.Unmarshal(reportBytes, &gitleaks)
	if err != nil {
		fmt.Println("Error decoding gitleaks report:", err)
	}

	// Print the parsed JSON report
	// fmt.Printf("Gitleaks JSON report for repository %s:\n", repoName)
	// for _, leak := range gitleaks {
	// 	fmt.Printf("Description: %s\n", leak.Description)
	// 	fmt.Printf("File: %s\nLine: %d\nOffender: %s\n\n", leak.File, leak.StartLine, leak.Secret)
	// }
	// fmt.Printf("Gitleaks output for repository %s:\n%s\n", repoName, output)

	fmt.Printf("Cleaning up repository %s...\n", repoName)
	cmd = exec.Command("rm", "-rf", repoName)
	cmd.Dir = "./repos"
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error cleaning up repository:", err)
	}
	return gitleaks, nil
}
