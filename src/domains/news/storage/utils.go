package storage

import (
	"gfeed/domains/news"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func toKeyAttributes(entry news.Entry) (map[string]*dynamodb.AttributeValue, error) {
	hash, err := entry.Hash()
	return map[string]*dynamodb.AttributeValue{
		"type": {
			S: aws.String(entry.Type),
		},
		"hash": {
			S: aws.String(hash),
		},
	}, err
}
