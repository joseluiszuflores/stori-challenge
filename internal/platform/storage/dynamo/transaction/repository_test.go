//go:build integration
// +build integration

package transaction

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang/glog"
	"github.com/joseluiszuflores/stori-challenge/internal"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/conn"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRepository_SaveTransaction(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		ctx         context.Context
		transaction internal.Transaction
	}
	conf, err := conn.NewAWSConfig("", "", "us-east-2", "http://localhost:8000", true)
	if err != nil {
		assert.NoError(t, err)

		return
	}
	today := time.Now()
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Success saving transaction",
			fields: fields{client: conn.NewDynamoDBClient(conf)},
			args: args{
				ctx: context.TODO(),
				transaction: internal.Transaction{
					ID:          2,
					Date:        today,
					Transaction: 900,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				client: tt.fields.client,
			}
			if err := r.SaveTransaction(tt.args.ctx, tt.args.transaction); (err != nil) != tt.wantErr {
				t.Errorf("SaveTransaction() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			data := make(map[string]types.AttributeValue)
			data["id"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", tt.args.transaction.ID)}
			data["date"] = &types.AttributeValueMemberS{Value: today.String()}
			//nolint:exhaustruct
			input := &dynamodb.GetItemInput{
				Key:       data,
				TableName: aws.String(tableNameTransactionDynamo),
			}

			out, err := tt.fields.client.GetItem(tt.args.ctx, input)
			if err != nil {
				assert.NoError(t, err)

				return
			}
			glog.Info(out)
			_ = out

		})
	}
}

func TestRepository_SaveTransactions(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	conf, err := conn.NewAWSConfig("", "", "us-east-1", "http://localhost:8000", true)
	if err != nil {
		assert.NoError(t, err)

		return
	}
	transactionsArr := make([]internal.Transaction, 0)
	transactionsArr = append(transactionsArr, internal.Transaction{
		ID:          1,
		Date:        time.Now(),
		Transaction: 90,
	})
	transactionsArr = append(transactionsArr, internal.Transaction{
		ID:          2,
		Date:        time.Now(),
		Transaction: 10,
	})
	transactionsArr = append(transactionsArr, internal.Transaction{
		ID:          3,
		Date:        time.Now(),
		Transaction: -11,
	})

	type args struct {
		ctx          context.Context
		transactions internal.Transactions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success saving all data in transactions table",
			fields: fields{
				client: conn.NewDynamoDBClient(conf),
			},
			args: args{
				ctx:          context.TODO(),
				transactions: transactionsArr,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				client: tt.fields.client,
			}
			tt.wantErr(t, r.SaveTransactions(tt.args.ctx, tt.args.transactions), fmt.Sprintf("SaveTransactions(%v, %v)", tt.args.ctx, tt.args.transactions))
			helperFnFindTransactioninDB(tt.args.ctx, t, tt.fields.client, tt.args.transactions)
		})
	}
}

func helperFnFindTransactioninDB(ctx context.Context, t *testing.T, client *dynamodb.Client, tns internal.Transactions) {
	t.Helper()
	for _, val := range tns {
		data := make(map[string]types.AttributeValue)
		data["id"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", val.ID)}
		data["date"] = &types.AttributeValueMemberS{Value: val.Date.String()}
		//nolint:exhaustruct
		input := &dynamodb.GetItemInput{
			Key:       data,
			TableName: aws.String(tableNameTransactionDynamo),
		}

		out, err := client.GetItem(ctx, input)
		if err != nil {
			assert.NoError(t, err)

			return
		}
		glog.Info(out)
	}
}
