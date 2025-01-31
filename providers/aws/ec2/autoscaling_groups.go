package ec2

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	. "github.com/tailwarden/komiser/models"
	. "github.com/tailwarden/komiser/providers"
)

const AWS_SERVICE_NAME_ASG = "AutoScalingGroup"

type AutoScalingGroupClient interface {
	DescribeAutoScalingGroups(
		ctx context.Context,
		params *autoscaling.DescribeAutoScalingGroupsInput,
		optFns ...func(*autoscaling.Options),
	) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}

func AutoScalingGroups(ctx context.Context, clientProvider ProviderClient) ([]Resource, error) {
	client := autoscaling.NewFromConfig(*clientProvider.AWSClient)

	d := ASGDiscoverer{
		Client:      client,
		Ctx:         ctx,
		AccountName: clientProvider.Name,
		Region:      clientProvider.AWSClient.Region,
	}

	return d.Discover()
}

type ASGDiscoverer struct {
	Client      AutoScalingGroupClient
	Ctx         context.Context
	AccountName string
	Region      string
}

func (d ASGDiscoverer) Discover() ([]Resource, error) {
	resources := make([]Resource, 0)
	var queryInput autoscaling.DescribeAutoScalingGroupsInput

	for {
		output, err := d.Client.DescribeAutoScalingGroups(d.Ctx, &queryInput)
		if err != nil {
			return resources, err
		}

		for _, asg := range output.AutoScalingGroups {
			tags := make([]Tag, 0)
			for _, tag := range asg.Tags {
				tags = append(tags, Tag{
					Key:   *tag.Key,
					Value: *tag.Value,
				})
			}

			jsonData, err := json.Marshal(asg)
			if err != nil {
				log.Printf("ERROR: Failed to marshall json: %v", err)
			}
			jsonString := string(jsonData)

			resources = append(resources, Resource{
				Provider:   "AWS",
				Account:    d.AccountName,
				Service:    AWS_SERVICE_NAME_ASG,
				Region:     d.Region,
				ResourceId: *asg.AutoScalingGroupARN,
				Cost:       0,
				Name:       *asg.AutoScalingGroupName,
				FetchedAt:  time.Now(),
				Tags:       tags,
				Data:       jsonString,
				Link: fmt.Sprintf(
					"https://%s.console.aws.amazon.com/ec2/home?region=%s#AutoScalingGroupDetails:id=%s",
					d.Region,
					d.Region,
					*asg.AutoScalingGroupName,
				),
			})
		}
		if aws.ToString(output.NextToken) == "" {
			break
		}

		queryInput.NextToken = output.NextToken
	}

	log.WithFields(log.Fields{
		"provider":  "AWS",
		"account":   d.AccountName,
		"region":    d.Region,
		"service":   "AutoScalingGroup",
		"resources": len(resources),
	}).Info("Fetched resources")
	return resources, nil
}
