package internal

type Balance struct {
	Total               float64
	AverageDebitAmount  float64
	AverageCreditAmount float64
	TransactionByMonth  map[string]int
}
