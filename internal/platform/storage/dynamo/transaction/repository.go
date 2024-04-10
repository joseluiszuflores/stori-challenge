package transaction

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
)

const tableNameTransactionDynamo = "transaction"

type Repository struct {
	client *dynamodb.Client
}

// SaveTransaction saves the data of transaction.
func (r *Repository) SaveTransaction(ctx context.Context, transaction mooc.Transaction) error {
	data := make(map[string]types.AttributeValue)
	data["id"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", transaction.ID)}
	data["date"] = &types.AttributeValueMemberS{Value: transaction.Date.String()}
	data["transaction"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", transaction.Transaction)}
	input := &dynamodb.PutItemInput{
		Item:      data,
		TableName: aws.String(tableNameTransactionDynamo),
	}

	_, err := r.client.PutItem(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// SaveTransactions saves the transactions into a transaction DB (TX).
func (r *Repository) SaveTransactions(ctx context.Context, transaction mooc.Transactions) error {
	transactionDB := make([]types.TransactWriteItem, 0)
	for _, trans := range transaction {
		data := make(map[string]types.AttributeValue)
		data["id"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", trans.ID)}
		data["date"] = &types.AttributeValueMemberS{Value: trans.Date.String()}
		data["transaction"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", trans.Transaction)}
		transactionDB = append(transactionDB, types.TransactWriteItem{
			Put: &types.Put{
				Item:      data,
				TableName: aws.String(tableNameTransactionDynamo),
			},
		})
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactionDB,
	}

	_, err := r.client.TransactWriteItems(ctx, input)
	if err != nil {
		return err
	}
	return nil
}
