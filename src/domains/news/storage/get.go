package storage

import (
	"gfeed/domains/news"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Get entry from DynamoDB
func Get(entry news.Entry) (*dynamodb.GetItemOutput, error) {
	key, err := toKeyAttributes(entry)

	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String(dynamoDBTable),
	}

	return svc.GetItem(input)
}

// IsRecorded checks if entry exist in database
func IsRecorded(e news.Entry) (bool, error) {
	r, err := Get(e)

	return len(r.Item) > 0, err
}
