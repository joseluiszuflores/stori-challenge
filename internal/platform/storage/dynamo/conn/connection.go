package conn

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//nolint:lll,revive
func NewAWSConfig(awsAccessKey, awsSecretAccessKey, awsRegion, urlDevAWSDynamo string, devEnv bool) (aws.Config, error) {

	if devEnv {
		return devAWSConf(awsRegion, urlDevAWSDynamo)
	}

	conf, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		return aws.Config{}, err
	}

	return conf, nil
}

func NewDynamoDBClient(sdkConfig aws.Config) *dynamodb.Client {
	// initialize new dynamodb client from lambda config and return it
	return dynamodb.NewFromConfig(sdkConfig)
}

func devAWSConf(awsRegion, urlDevAWSDynamo string) (aws.Config, error) {
	// Create a custom endpoint resolver function
	endpointResolver := aws.EndpointResolverWithOptionsFunc(
		//nolint: revive
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			//nolint:exhaustruct
			return aws.Endpoint{
				URL:           urlDevAWSDynamo,
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
