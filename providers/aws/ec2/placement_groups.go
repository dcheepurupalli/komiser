package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"
)

func PlacementGroups(ctx context.Context, client providers.ProviderClient) ([]models.Resource, error) {
	resources := make([]models.Resource, 0)
	ec2Client := ec2.NewFromConfig(*client.AWSClient)

	input := &ec2.DescribePlacementGroupsInput{}
	placementGroups, err := ec2Client.DescribePlacementGroups(ctx, input)
	if err != nil {
		return resources, err
	}

	for _, placementGroup := range placementGroups.PlacementGroups {

		var tags []models.Tag
		for _, tag := range placementGroup.Tags {
			tags = append(tags, models.Tag{
				Key:   aws.ToString(tag.Key),
				Value: aws.ToString(tag.Value),
			})
		}

		jsonData, err := json.Marshal(placementGroup)
		if err != nil {
			log.Printf("ERROR: Failed to marshall json: %v", err)
		}
		jsonString := string(jsonData)

		resources = append(resources, models.Resource{
			Provider:   "AWS",
			Account:    client.Name,
			Service:    "EC2 Placement Group",
			ResourceId: aws.ToString(placementGroup.GroupName),
			Region:     client.AWSClient.Region,
			Name:       aws.ToString(placementGroup.GroupName),
			Tags:       tags,
			Data:       jsonString,
			FetchedAt:  time.Now(),
			Link:       fmt.Sprintf("https://%s.console.aws.amazon.com/ec2/v2/home?region=%s#PlacementGroups:sort=groupName", client.AWSClient.Region, client.AWSClient.Region),
		})
	}

	log.WithFields(log.Fields{
		"provider":  "AWS",
		"account":   client.Name,
		"region":    client.AWSClient.Region,
		"service":   "EC2 Placement Group",
		"resources": len(resources),
	}).Info("Fetched resources")

	return resources, nil
}
