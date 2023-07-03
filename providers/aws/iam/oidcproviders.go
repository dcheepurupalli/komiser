package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	. "github.com/tailwarden/komiser/models"
	. "github.com/tailwarden/komiser/providers"
)

func OIDCProviders(ctx context.Context, client ProviderClient) ([]Resource, error) {
	resources := make([]Resource, 0)
	var config iam.ListOpenIDConnectProvidersInput
	iamClient := iam.NewFromConfig(*client.AWSClient)

	output, err := iamClient.ListOpenIDConnectProviders(ctx, &config)
	if err != nil {
		return resources, err
	}

	for _, oidcprovider := range output.OpenIDConnectProviderList {
		outputProvider, err := iamClient.GetOpenIDConnectProvider(ctx, &iam.GetOpenIDConnectProviderInput{
			OpenIDConnectProviderArn: oidcprovider.Arn,
		})
		if err != nil {
			return resources, err
		}

		tags := make([]Tag, 0)
		for _, t := range outputProvider.Tags {
			tags = append(tags, Tag{
				Key:   *t.Key,
				Value: *t.Value,
			})
		}

		jsonData, err := json.Marshal(oidcprovider)
		if err != nil {
			log.Printf("ERROR: Failed to marshall json: %v", err)
		}
		jsonString := string(jsonData)

		resources = append(resources, Resource{
			Provider:   "AWS",
			Account:    client.Name,
			Service:    "IAM Identity Provider",
			ResourceId: *oidcprovider.Arn,
			Region:     client.AWSClient.Region,
			Name:       *outputProvider.Url,
			Cost:       0,
			Tags:       tags,
			Data:       jsonString,
			CreatedAt:  *outputProvider.CreateDate,
			FetchedAt:  time.Now(),
			Link:       fmt.Sprintf("https://%s.console.aws.amazon.com/iamv2/home?region=%s#/identity_providers/details/%s", client.AWSClient.Region, client.AWSClient.Region, *oidcprovider.Arn),
		})
	}

	log.WithFields(log.Fields{
		"provider":  "AWS",
		"account":   client.Name,
		"region":    client.AWSClient.Region,
		"service":   "IAM Identity Provider",
		"resources": len(resources),
	}).Info("Fetched resources")

	return resources, nil
}
