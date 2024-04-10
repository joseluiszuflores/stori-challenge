package internal

import (
	"context"
	"time"
)

type Balance struct {
	Total               float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	TransactionByMonth  map[string]int
}

type Transaction struct {
	ID          int
	Date        time.Time
	Transaction float64
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Transactions []Transaction

func (t Transactions) Sum() float64 {
	sum := 0.0
	for _, val := range t {
		sum += val.Transaction
	}

	return sum
}

//go:generate mockgen -destination=./mocks/transaction_mock.go -package=mocks -source=./balance.go RepositoryTransaction

type RepositoryTransaction interface {
	SaveTransaction(ctx context.Context, transaction Transaction) error
}

//go:generate mockgen -destination=./mocks/email_mock.go -package=mocks -source=./balance.go EmailService
type EmailService interface {
	Send(destination, name string, balance Balance) error
}
