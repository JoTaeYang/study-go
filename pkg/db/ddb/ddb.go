package ddb

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ratelimit"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Config struct {
	Region    string
	AccessKey string
	SecretKey string
	Table     string
}

var (
	DBTableName  string
	Dynamodbconn *dynamodb.Client
)

func LocalInit(cfg Config) error {

	endpointResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: "http://localhost:4566",
			//SigningRegion: region,
		}, nil
	})

	dbCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(endpointResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "4nbu8j", SecretAccessKey: "sdglg",
			},
		}),
		config.WithRetryer(func() aws.Retryer {
			return retry.NewStandard(func(so *retry.StandardOptions) {
				so.RateLimiter = ratelimit.NewTokenRateLimit(1000000)
			})
		}),
	)
	if err != nil {
		log.Println(err)
		return err
	}
	Dynamodbconn = dynamodb.NewFromConfig(dbCfg)
	if Dynamodbconn == nil {
		return errors.New("dynamodb connection fail")
	}

	log.Println("local dynamodb init")
	return nil
}

func Init(cfg Config) error {
	if cfg.AccessKey != "" {
		awsConfig, err := config.LoadDefaultConfig(
			context.Background(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
			config.WithRegion(cfg.Region),
		)
		if err != nil {
			log.Println(err)
			return err
		}
		Dynamodbconn = dynamodb.NewFromConfig(awsConfig)
	}

	awsConfig, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		log.Println(err)
		return err
	}

	Dynamodbconn = dynamodb.NewFromConfig(awsConfig)

	if Dynamodbconn == nil {
		return errors.New("dynamodb connection fail")
	}

	log.Println("dynamodb init")
	return nil
}
