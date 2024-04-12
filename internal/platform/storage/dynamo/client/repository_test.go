package client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
	"github.com/joseluiszuflores/stori-challenge/internal/platform/storage/dynamo/conn"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRepository_GetClient(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		ctx context.Context
		id  int
	}
	conf, err := conn.NewAWSConfig("", "", "us-east-2", "http://localhost:8000", true)
	if err != nil {
		assert.NoError(t, err)

		return
	}

	db := conn.NewDynamoDBClient(conf)

	usr := createUser(t, db)
	if usr == nil {
		assert.NotNil(t, usr)

		return
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *mooc.User
		wantErr bool
	}{
		{
			name:   "Success get the user",
			fields: fields{client: db},
			args: args{
				ctx: context.TODO(),
				id:  usr.ID,
			},
			want:    usr,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				client: tt.fields.client,
			}
			got, err := r.GetClient(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createUser(t *testing.T, db *dynamodb.Client) *mooc.User {
	t.Helper()
	usr := mooc.User{
		ID:    1,
		Name:  "Jose Luis",
		Email: "storymockexample@gmail.com",
	}

	data := make(map[string]types.AttributeValue)
	data["id"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", usr.ID)}
	data["email"] = &types.AttributeValueMemberS{Value: usr.Email}
	data["name"] = &types.AttributeValueMemberS{Value: usr.Name}
	input := &dynamodb.PutItemInput{
		Item:      data,
		TableName: aws.String(tableNameUserDynamo),
	}
	_, err := db.PutItem(context.TODO(), input)
	assert.NoError(t, err)
	return &usr
}
