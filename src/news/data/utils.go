package data

import (
	"gfeed/news"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func toKeyAttributes(entry news.Entry) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"type": {
			S: aws.String(entry.Type),
		},
		"hash": {
			S: aws.String(entry.Hash()),
		},
	}
}
