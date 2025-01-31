package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	. "github.com/tailwarden/komiser/models"
	. "github.com/tailwarden/komiser/providers"
)

func ElasticIps(ctx context.Context, client ProviderClient) ([]Resource, error) {
	config := ec2.DescribeAddressesInput{}
	resources := make([]Resource, 0)
	ec2Client := ec2.NewFromConfig(*client.AWSClient)
	configClient := configservice.NewFromConfig(*client.AWSClient)

	stsClient := sts.NewFromConfig(*client.AWSClient)
	stsOutput, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return resources, err
	}

	accountId := stsOutput.Account

	for {
		output, err := ec2Client.DescribeAddresses(ctx, &config)
		if err != nil {
			return resources, err
		}

		for _, elasticIps := range output.Addresses {
			tags := make([]Tag, 0)
			for _, tag := range elasticIps.Tags {
				tags = append(tags, Tag{
					Key:   *tag.Key,
					Value: *tag.Value,
				})
			}

			cost := 0.0

			resourceConfig, err := configClient.BatchGetResourceConfig(ctx, &configservice.BatchGetResourceConfigInput{
				ResourceKeys: []types.ResourceKey{
					{
						ResourceId:   elasticIps.AllocationId,
						ResourceType: "AWS::EC2::EIP",
					},
				},
			})
			if err != nil {
				log.WithFields(log.Fields{
					"service":  "Elastic IP",
					"name":     *elasticIps.AllocationId,
					"region":   client.AWSClient.Region,
					"provider": "AWS",
				}).Warn("Cost couldn't be calculated due to missing AWS config")
			} else {
				if len(resourceConfig.BaseConfigurationItems) > 0 {
					creationTime := resourceConfig.BaseConfigurationItems[0].ResourceCreationTime
					if creationTime != nil {
						hoursSinceCreation := hoursSince(*creationTime)

						hourlyCost := 0.005

						if elasticIps.InstanceId != nil {
							cost = 0
						} else {
							cost = hourlyCost * hoursSinceCreation
						}
					} else {
						cost = 0
						log.WithFields(log.Fields{
							"service":  "Elastic IP",
							"provider": "AWS",
						}).Error("Cost couldn't be calculated because the creationTime returned by resource config is nil")
					}

				}
			}

			jsonData, err := json.Marshal(elasticIps)
			if err != nil {
				log.Printf("ERROR: Failed to marshall json: %v", err)
			}
			jsonString := string(jsonData)

			resourceArn := fmt.Sprintf("arn:aws:ec2:%s:%s:elastic-ip/%s", client.AWSClient.Region, *accountId, *elasticIps.AllocationId)

			resources = append(resources, Resource{
				Provider:   "AWS",
				Account:    client.Name,
				Service:    "Elastic IP",
				Region:     client.AWSClient.Region,
				ResourceId: resourceArn,
				Cost:       cost,
				Name:       *elasticIps.AllocationId,
				FetchedAt:  time.Now(),
				Tags:       tags,
				Data:       jsonString,
				Link:       fmt.Sprintf("https:/%s.console.aws.amazon.com/ec2/home?region=%s#ElasticIpDetails:AllocationId=%s", client.AWSClient.Region, client.AWSClient.Region, *elasticIps.AllocationId),
			})
		}

		log.WithFields(log.Fields{
			"provider":  "AWS",
			"account":   client.Name,
			"region":    client.AWSClient.Region,
			"service":   "Elastic IP",
			"resources": len(resources),
		}).Info("Fetched resources")
		return resources, nil
	}
}

func hoursSince(t time.Time) float64 {
	duration := time.Since(t)
	return duration.Hours()
}
