package data

import (
	"gfeed/news"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const ttlTime = time.Hour * 360 // 15 days

// Put news on table
func Put(entry news.Entry) error {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	ttl := strconv.FormatInt(time.Now().Add(ttlTime).Unix(), 10)

	item := toKeyAttributes(entry)

	item["title"] = &dynamodb.AttributeValue{
		S: aws.String(entry.Title),
	}

	item["date"] = &dynamodb.AttributeValue{
		N: aws.String(now),
	}

	item["ttl"] = &dynamodb.AttributeValue{
		N: aws.String(ttl),
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(dynamoDBTable),
	}

	_, err := svc.PutItem(input)

	return err
}
