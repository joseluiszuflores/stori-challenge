package migration

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/aws/smithy-go/transport/http"
	"github.com/golang/glog"
)

var (
	ErrCreatingClientTable      = errors.New("error to create client table")
	ErrCreatingTransactionTable = errors.New("error to create transaction table")
)

// Migrator will create and  check if the tablas for the system are exist.
// in case that these are not exist  migrator will send the instructions to  create it.
type Migrator struct {
	client *dynamodb.Client
}

func NewMigrator(client *dynamodb.Client) *Migrator {
	return &Migrator{client: client}
}

func (m *Migrator) Do() error {
	glog.Info("Doing migration")
	if out, err := m.CreateClientTable(); err != nil || out == nil {
		glog.Error(err)

		return ErrCreatingClientTable
	}
	if out, err := m.CreateTransactionTable(); err != nil || out == nil {
		glog.Error(err)

		return ErrCreatingClientTable
	}

	return nil
}

func (m *Migrator) DescribeTable(tableName string) (*dynamodb.DescribeTableOutput, error) {
	// get table information
	table, err := m.client.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
	)

	// return table and error
	return table, err
}

const (
	tableClients     = "user"
	tableTransaction = "transaction"
)

// you need to add any other attributes to items, you add those on item insert, not table creation.
// DynamoDB only really cares about and enforces the primary key.
func (m *Migrator) CreateClientTable() (*dynamodb.CreateTableOutput, error) {
	des, err := m.DescribeTable(tableClients)
	if err != nil && !validationErrorNotFoundTable(err) {
		return nil, err
	}
	// the table should be to exist.
	if des != nil {
		//nolint: exhaustruct
		return &dynamodb.CreateTableOutput{}, nil
	}
	// create client table
	//nolint: exhaustruct
	table, err := m.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableClients),
		// primary key attributes are required
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		// add primary key details
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		// set billing mode
		// Using on-demand provisioning (pay as you go, no pre-allocation).
		BillingMode: types.BillingModePayPerRequest,
	})

	return table, err
}

func validationErrorNotFoundTable(err error) bool {
	var operationError *smithy.OperationError
	if !errors.As(err, &operationError) {
		return false
	}
	var responseError *http.ResponseError
	if !errors.As(operationError.Err, &responseError) {
		return false
	}
	var notfound *types.ResourceNotFoundException
	if !errors.As(responseError.Err, &notfound) {
		return false
	}
	if notfound != nil {
		return true
	}

	return false
}

func (m *Migrator) CreateTransactionTable() (*dynamodb.CreateTableOutput, error) {
	des, err := m.DescribeTable(tableTransaction)
	if err != nil && !validationErrorNotFoundTable(err) {
		return nil, err
	}
	if des != nil {
		//nolint:exhaustruct
		return &dynamodb.CreateTableOutput{}, nil
	}
	// create client table.
	//nolint: exhaustruct
	table, err := m.client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableTransaction),
		// primary key attributes are required.
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		// add primary key details.
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			}, {
				AttributeName: aws.String("date"),
				KeyType:       types.KeyTypeRange,
			},
		},
		// set billing mode.
		// Using on-demand provisioning (pay as you go, no pre-allocation).
		BillingMode: types.BillingModePayPerRequest,
	})

	return table, err
}
