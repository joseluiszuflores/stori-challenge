package internal

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

type Balance struct {
	Total               float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	TransactionByMonth  map[string]int
	AverageByMonth      map[string]*Average
}

type Average struct {
	AverageDebitAmount  float64
	AverageCreditAmount float64
}

//nolint:lll
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

//go:generate mockgen -destination=./mocks/mock.go -package=mocks -source=./balance.go RepositoryTransaction

type RepositoryTransaction interface {
	SaveTransaction(ctx context.Context, transaction Transaction) error
	SaveTransactions(ctx context.Context, transaction Transactions) error
}

//go:generate mockgen -destination=./mocks/mock.go -package=mocks -source=./balance.go RepositoryUser

type RepositoryUser interface {
	GetClient(ctx context.Context, id int) (*User, error)
}

//go:generate mockgen -destination=./mocks/mock.go -package=mocks -source=./balance.go EmailService

type EmailService interface {
	Send(destination, name string, balance Balance) error
}

func ToInt(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return -1
	}

	return intVal
}

func IsValidInt(val int) bool {
	return val == -1
}

var ErrIDUserIsInvalid = errors.New("the user id is invalid")

func ToIntFromFile(val string) string {
	allName := strings.Split(val, ".")
	if len(allName) < 1 {
		return ""
	}
	return allName[0]
}
