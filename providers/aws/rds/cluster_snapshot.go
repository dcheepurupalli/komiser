package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/tailwarden/komiser/models"
	"github.com/tailwarden/komiser/providers"
	"github.com/tailwarden/komiser/utils"
)

func ClusterSnapshots(ctx context.Context, client providers.ProviderClient) ([]models.Resource, error) {
	var config rds.DescribeDBClusterSnapshotsInput
	resources := make([]models.Resource, 0)
	rdsClient := rds.NewFromConfig(*client.AWSClient)

	for {
		output, err := rdsClient.DescribeDBClusterSnapshots(ctx, &config)
		if err != nil {
			return resources, err
		}

		for _, clusterSnapshot := range output.DBClusterSnapshots {
			tags := make([]models.Tag, 0)
			for _, tag := range clusterSnapshot.TagList {
				tags = append(tags, models.Tag{
					Key:   *tag.Key,
					Value: *tag.Value,
				})
			}
			_clusterSnapshotName := *clusterSnapshot.DBClusterSnapshotIdentifier

			startOfMonth := utils.BeginningOfMonth(time.Now())
			hourlyUsage := 0
			if (*clusterSnapshot.ClusterCreateTime).Before(startOfMonth) {
				hourlyUsage = int(time.Since(startOfMonth).Hours())
			} else {
				hourlyUsage = int(time.Since(*clusterSnapshot.ClusterCreateTime).Hours())
			}

			hourlyCost := 0.0
			monthlyCost := float64(hourlyUsage) * hourlyCost

			jsonData, err := json.Marshal(clusterSnapshot)
			if err != nil {
				log.Printf("ERROR: Failed to marshall json: %v", err)
			}
			jsonString := string(jsonData)

			resources = append(resources, models.Resource{
				Provider:   "AWS",
				Account:    client.Name,
				Service:    "RDS Cluster Snapshot",
				Region:     client.AWSClient.Region,
				ResourceId: *clusterSnapshot.DBClusterSnapshotArn,
				Cost:       monthlyCost,
				Name:       _clusterSnapshotName,
				FetchedAt:  time.Now(),
				Tags:       tags,
				Data:       jsonString,
				Link:       fmt.Sprintf("https:/%s.console.aws.amazon.com/rds/home?region=%s#snapshots-list:id=%s", client.AWSClient.Region, client.AWSClient.Region, *clusterSnapshot.DBClusterSnapshotIdentifier),
			})
		}

		if aws.ToString(output.Marker) == "" {
			break
		}

		config.Marker = output.Marker
	}
	log.WithFields(log.Fields{
		"provider":  "AWS",
		"account":   client.Name,
		"region":    client.AWSClient.Region,
		"service":   "RDS Cluster Snapshot",
		"resources": len(resources),
	}).Info("Fetched resources")
	return resources, nil
}
