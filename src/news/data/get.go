package data

import (
	"gfeed/news"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const dynamoDBTable = "gamer-feed-test"

// Get entry from DynamoDB
func Get(entry news.Entry) (*dynamodb.GetItemOutput, error) {
	input := &dynamodb.GetItemInput{
		Key:       toKeyAttributes(entry),
		TableName: aws.String(dynamoDBTable),
	}

	return svc.GetItem(input)
}

// IsRecorded checks if entry exist in database
func IsRecorded(e news.Entry) (bool, error) {
	r, err := Get(e)

	return len(r.Item) > 0, err
}
