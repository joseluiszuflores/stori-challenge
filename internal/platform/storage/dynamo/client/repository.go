package client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
)

const tableNameUserDynamo = "user"

type Repository struct {
	client *dynamodb.Client
}

func NewRepository(client *dynamodb.Client) *Repository {
	return &Repository{client: client}
}

func (r *Repository) GetClient(ctx context.Context, id int) (*mooc.User, error) {
	data := make(map[string]types.AttributeValue)
	data["id"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", id)}

	input := &dynamodb.GetItemInput{
		Key:       data,
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
