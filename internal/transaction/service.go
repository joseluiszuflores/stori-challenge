package transaction

import (
	"context"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
)

type Service struct {
	idUser       int
	transactions mooc.Transactions
	debit        mooc.Transactions
	credit       mooc.Transactions
	email        mooc.EmailService
	userRep      mooc.RepositoryUser
	transRep     mooc.RepositoryTransaction
}

func NewService(idUser int, transactions mooc.Transactions, debit mooc.Transactions, credit mooc.Transactions, email mooc.EmailService, userRep mooc.RepositoryUser, transRep mooc.RepositoryTransaction) *Service {
	return &Service{idUser: idUser, transactions: transactions, debit: debit, credit: credit, email: email, userRep: userRep, transRep: transRep}
}

func (s *Service) SummaryTransaction(ctx context.Context) error {
	usr, err := s.userRep.GetClient(ctx, s.idUser)
	if err != nil {
		return err
	}

	s.SeparatedDebitCredit()
	balance := mooc.Balance{
		Total:               s.transactions.Sum(),
		AverageDebitAmount:  s.AverageDebit(),
		AverageCreditAmount: s.AverageCredit(),
		TransactionByMonth:  s.MovementsByMonth(),
	}

	if err := s.email.Send(usr.Email, usr.Name, balance); err != nil {
		return err
	}

	if err := s.transRep.SaveTransactions(ctx, s.transactions); err != nil {
		return err
	}

	return nil
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

func (s *Service) MovementsByMonth() map[string]int {
	months := make(map[string]int)
	for _, transaction := range s.transactions {
		_, ok := months[mooc.Months[int(transaction.Date.Month())]]
		if ok {
			months[mooc.Months[int(transaction.Date.Month())]]++
		} else {
			months[mooc.Months[int(transaction.Date.Month())]] = 1
		}
	}

	return months
}
