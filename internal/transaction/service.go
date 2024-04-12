package transaction

import (
	"context"
	"github.com/golang/glog"
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

//nolint:lll
func NewService(idUser string, transactions mooc.Transactions, email mooc.EmailService, userRep mooc.RepositoryUser, transRep mooc.RepositoryTransaction) (*Service, error) {
	id := mooc.ToInt(idUser)
	if mooc.IsValidInt(id) {
		return nil, mooc.ErrIDUserIsInvalid
	}
	return &Service{idUser: id, transactions: transactions, email: email, userRep: userRep, transRep: transRep}, nil
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
		AverageByMonth:      s.SumAverageByMonth(s.SeparateTransactionByMonth()),
	}
	glog.Info("Sending email")
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

type MountsByMonth struct {
	AverageDebitAmount  mounts
	AverageCreditAmount mounts
}

type mounts []float64

func (m mounts) Sum() float64 {
	sum := 0.0
	for _, val := range m {
		sum += val
	}

	return sum
}

func (s *Service) SeparateTransactionByMonth() map[string]*MountsByMonth {
	months := make(map[string]*MountsByMonth)
	for _, transaction := range s.transactions {
		val, ok := months[mooc.Months[int(transaction.Date.Month())]]
		if !ok {
			months[mooc.Months[int(transaction.Date.Month())]] = &MountsByMonth{
				AverageCreditAmount: make(mounts, 0),
				AverageDebitAmount:  make(mounts, 0),
			}
			val = months[mooc.Months[int(transaction.Date.Month())]]
		}
		if transaction.Transaction < 0 {
			val.AverageDebitAmount = append(val.AverageDebitAmount, transaction.Transaction)
		} else {
			val.AverageCreditAmount = append(val.AverageCreditAmount, transaction.Transaction)
		}
	}

	return months
}

func (s *Service) SumAverageByMonth(mountsBymonth map[string]*MountsByMonth) map[string]*mooc.Average {
	months := make(map[string]*mooc.Average)
	for month, mounts := range mountsBymonth {
		monthVal, ok := months[month]
		if !ok {
			months[month] = &mooc.Average{
				AverageDebitAmount:  0,
				AverageCreditAmount: 0,
			}
			monthVal = months[month]
		}
		monthVal.AverageDebitAmount = mounts.AverageDebitAmount.Sum() / float64(len(mounts.AverageDebitAmount))
		monthVal.AverageCreditAmount = mounts.AverageCreditAmount.Sum() / float64(len(mounts.AverageCreditAmount))
	}

	return months
}
