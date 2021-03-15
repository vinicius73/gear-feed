package storage

import (
	"gfeed/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var svc *dynamodb.DynamoDB
var dynamoDBTable string

func init() {
	region := utils.GetEnv("AWS_DEFAULT_REGION", "us-east-1")

	sess := session.Must(
		session.NewSessionWithOptions(
			session.Options{
				Config: aws.Config{Region: aws.String(region)},
			},
		),
	)

	svc = dynamodb.New(sess)

	dynamoDBTable = utils.GetEnv("ENTRIES_TABLE", "gamer-feed-entries")
}
