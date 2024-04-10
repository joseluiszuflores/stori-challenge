package client

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
)

const tableNameUserDynamo = "user"

type Repository struct {
	client *dynamodb.Client
}

type userDTO struct {
	ID uuid.UUID
}

func (r *Repository) GetClient(ctx context.Context, id uuid.UUID) (*mooc.User, error) {
	userToGet, err := attributevalue.MarshalMap(userDTO{ID: id})
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		Key:       userToGet,
		TableName: aws.String(tableNameUserDynamo),
	}
	item, err := r.client.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	var usr mooc.User
	if err := attributevalue.UnmarshalMap(item.Item, &usr); err != nil {
		return nil, err
	}
	return &usr, nil
}
