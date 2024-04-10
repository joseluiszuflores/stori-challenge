package internal

import "time"

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

type Transactions []Transaction
