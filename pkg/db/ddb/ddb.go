package ddb

import (
	"context"
	"errors"
	"log"

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
	return nil
}
