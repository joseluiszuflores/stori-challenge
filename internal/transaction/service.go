package transaction

import mooc "github.com/joseluiszuflores/stori-challenge/internal"

type Service struct {
	transactions mooc.Transactions
	debit        mooc.Transactions
	credit       mooc.Transactions
}

func (s *Service) SeparatedDebitCredit() {
	for _, transaction := range s.transactions {
		if transaction.Transaction < 0 {
			s.debit = append(s.debit, transaction)
		} else {
			s.credit = append(s.credit, transaction)
		}
	}
}

func (s *Service) AverageDebit() float64 {
	return s.debit.Sum() / float64(len(s.debit))
}

func (s *Service) AverageCredit() float64 {
	return s.credit.Sum() / float64(len(s.credit))
}

func (s *Service) MovementsByMonth() map[int]int {
	months := make(map[int]int)
	for _, trasaction := range s.transactions {
		_, ok := months[int(trasaction.Date.Month())]
		if ok {
			months[int(trasaction.Date.Month())]++
		} else {
			months[int(trasaction.Date.Month())] = 1
		}
	}

	return months
}
