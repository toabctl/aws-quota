package common

import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func AwsRegions() []string {
	region_names := make([]string, 0)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	
	client := ec2.NewFromConfig(cfg)

	params := &ec2.DescribeRegionsInput{
	}

	regions, err := client.DescribeRegions(context.TODO(), params)
	if err != nil {
		log.Fatalf("unable to get regions, %v", err)
	}
	for _, r := range (regions.Regions) {
		region_names = append(region_names, aws.ToString(r.RegionName))
	}
	return region_names
}
