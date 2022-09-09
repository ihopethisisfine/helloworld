package aws

import (
	"context"
	"log"
	"time"

	"github.com/ihopethisisfine/helloworld/internal/domain"
	"github.com/ihopethisisfine/helloworld/internal/pkg/storage"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var _ storage.UserStorer = UserStorage{}

type UserStorage struct {
	timeout time.Duration
	client  *dynamodb.DynamoDB
}

func NewUserStorage(session *session.Session, timeout time.Duration) UserStorage {
	return UserStorage{
		timeout: timeout,
		client:  dynamodb.New(session),
	}
}

func (u UserStorage) Insert(ctx context.Context, user storage.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Println(err)
		return domain.ErrInternal
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      item,
		ExpressionAttributeNames: map[string]*string{
			"#username": aws.String("username"),
		},
		ConditionExpression: aws.String("attribute_not_exists(#username)"),
	}

	if _, err := u.client.PutItemWithContext(ctx, input); err != nil {
		log.Println(err)

		if _, ok := err.(*dynamodb.ConditionalCheckFailedException); ok {
			return domain.ErrConflict
		}

		return domain.ErrInternal
	}

	return nil
}

func (u UserStorage) Find(ctx context.Context, username string) (storage.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	input := &dynamodb.GetItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {S: aws.String(username)},
		},
	}

	res, err := u.client.GetItemWithContext(ctx, input)
	if err != nil {
		log.Println(err)

		return storage.User{}, domain.ErrInternal
	}

	if res.Item == nil {
		return storage.User{}, domain.ErrNotFound
	}

	var user storage.User
	if err := dynamodbattribute.UnmarshalMap(res.Item, &user); err != nil {
		log.Println(err)

		return storage.User{}, domain.ErrInternal
	}

	return user, nil
}