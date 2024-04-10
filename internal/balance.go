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

var Months = map[int]string{1: "January", 2: "February", 3: "March", 4: "April", 5: "May", 6: "June", 7: "July", 8: "August", 9: "September", 10: "October", 11: "November", 12: "December"}

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

//go:generate mockgen -destination=./mocks/user_mock.go -package=mocks -source=./balance.go RepositoryUser

type RepositoryUser interface {
	GetClient(ctx context.Context, id int) (*User, error)
}

//go:generate mockgen -destination=./mocks/email_mock.go -package=mocks -source=./balance.go EmailService

type EmailService interface {
	Send(destination, name string, balance Balance) error
}
