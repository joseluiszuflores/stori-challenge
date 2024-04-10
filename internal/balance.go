package internal

import (
	"context"
	"github.com/google/uuid"
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
	ID    uuid.UUID
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

type RepositoryTransaction interface {
	SaveTransaction(ctx context.Context, transaction Transaction) error
}
