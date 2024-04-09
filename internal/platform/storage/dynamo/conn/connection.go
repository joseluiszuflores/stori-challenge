package conn

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewAWSConfig() (aws.Config, error) {
	// get config from environment variables
	awsAccessKey := "any"
	awsSecretAccessKey := "any"
	awsRegion := "us-east-1"
	// setup aws credential provider
	_ = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		awsAccessKey,
		awsSecretAccessKey,
		"",
	))
	// Create a custom endpoint resolver function
	endpointResolver := aws.EndpointResolverWithOptionsFunc(
		//nolint: revive
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			//nolint:exhaustruct
			return aws.Endpoint{
				URL:           "http://localhost:8000",
				SigningRegion: region,
			}, nil
		})

	conf, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(endpointResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			//nolint:exhaustruct
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return conf, nil
}

func NewDynamoDBClient(sdkConfig aws.Config) *dynamodb.Client {
	// initialize new dynamodb client from aws config and return it
	return dynamodb.NewFromConfig(sdkConfig)
}
