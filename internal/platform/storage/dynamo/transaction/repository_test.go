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
	conf, err := conn.NewAWSConfig()
	if err != nil {
		assert.NoError(t, err)

		return
	}
	d := time.Now()
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
					Date:        d,
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
			data["date"] = &types.AttributeValueMemberS{Value: d.String()}

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
