package transaction

import (
	"context"
	"fmt"
	"github.com/joseluiszuflores/stori-challenge/internal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_SeparatedDebitCredit(t *testing.T) {
	type fields struct {
		transactions internal.Transactions
		debit        internal.Transactions
		credit       internal.Transactions
	}

	debit := make(internal.Transactions, 0)
	timeForAny, err := time.Parse(time.DateOnly, "2024-07-24")
	if err != nil {
		assert.NoError(t, err)
		return
	}
	debit = append(debit, internal.Transaction{
		ID:          1,
		Date:        timeForAny,
		Transaction: -10,
	})
	debit = append(debit, internal.Transaction{
		ID:          2,
		Date:        timeForAny,
		Transaction: -30,
	})

	credit := make(internal.Transactions, 0)
	credit = append(credit, internal.Transaction{
		ID:          3,
		Date:        timeForAny,
		Transaction: 10,
	})
	credit = append(credit, internal.Transaction{
		ID:          3,
		Date:        timeForAny,
		Transaction: 30,
	})
	allTrans := make(internal.Transactions, 0)
	allTrans = append(allTrans, credit...)
	allTrans = append(allTrans, debit...)
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Success separated the debit and credit transactions",
			fields: fields{
				transactions: allTrans,
				debit:        nil,
				credit:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				transactions: tt.fields.transactions,
			}
			s.SeparatedDebitCredit()
			assert.Equal(t, float64(0), s.transactions.Sum())
			assert.Equal(t, -40.0, s.debit.Sum())
			assert.Equal(t, 40.0, s.credit.Sum())
		})
	}
}

func TestService_MovementsByMonth(t *testing.T) {
	type fields struct {
		transactions internal.Transactions
		debit        internal.Transactions
		credit       internal.Transactions
	}

	debit := make(internal.Transactions, 0)
	timeForDebit, err := time.Parse(time.DateOnly, "2024-07-24")
	if err != nil {
		assert.NoError(t, err)
		return
	}
	debit = append(debit, internal.Transaction{
		ID:          1,
		Date:        timeForDebit,
		Transaction: -10,
	})
	debit = append(debit, internal.Transaction{
		ID:          2,
		Date:        timeForDebit,
		Transaction: -30,
	})

	credit := make(internal.Transactions, 0)
	timeForCredit, err := time.Parse(time.DateOnly, "2024-06-24")
	if err != nil {
		assert.NoError(t, err)
		return
	}
	credit = append(credit, internal.Transaction{
		ID:          3,
		Date:        timeForCredit,
		Transaction: 10,
	})
	credit = append(credit, internal.Transaction{
		ID:          3,
		Date:        timeForCredit,
		Transaction: 30,
	})
	allTrans := make(internal.Transactions, 0)
	allTrans = append(allTrans, credit...)
	allTrans = append(allTrans, debit...)

	tests := []struct {
		name   string
		fields fields
		want   map[string]int
	}{
		{
			name: "Success get the number of movements by month",
			fields: fields{
				transactions: allTrans,
			},
			want: map[string]int{"July": 2, "June": 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				transactions: tt.fields.transactions,
				debit:        tt.fields.debit,
				credit:       tt.fields.credit,
			}
			assert.Equalf(t, tt.want, s.MovementsByMonth(), "MovementsByMonth()")
		})
	}
}

func TestService_AverageDebit(t *testing.T) {
	type fields struct {
		transactions internal.Transactions
		debit        internal.Transactions
		credit       internal.Transactions
	}

	debit := make(internal.Transactions, 0)
	timeForDebit, err := time.Parse(time.DateOnly, "2024-07-24")
	if err != nil {
		assert.NoError(t, err)
		return
	}
	debit = append(debit, internal.Transaction{
		ID:          1,
		Date:        timeForDebit,
		Transaction: -10.3,
	})
	debit = append(debit, internal.Transaction{
		ID:          2,
		Date:        timeForDebit,
		Transaction: -20.46,
	})

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "success get the average debit",
			fields: fields{
				debit: debit,
			},
			want: -15.38,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				transactions: tt.fields.transactions,
				debit:        tt.fields.debit,
				credit:       tt.fields.credit,
			}
			assert.Equalf(t, tt.want, s.AverageDebit(), "AverageDebit()")
		})
	}
}

func TestService_AverageCredit(t *testing.T) {
	type fields struct {
		transactions internal.Transactions
		debit        internal.Transactions
		credit       internal.Transactions
	}

	credit := make(internal.Transactions, 0)
	timeForCredit, err := time.Parse(time.DateOnly, "2024-06-24")
	if err != nil {
		assert.NoError(t, err)
		return
	}
	credit = append(credit, internal.Transaction{
		ID:          3,
		Date:        timeForCredit,
		Transaction: 60.5,
	})
	credit = append(credit, internal.Transaction{
		ID:          3,
		Date:        timeForCredit,
		Transaction: 10,
	})

	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Success get average Credit",
			fields: fields{
				credit: credit,
			},
			want: 35.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				transactions: tt.fields.transactions,
				debit:        tt.fields.debit,
				credit:       tt.fields.credit,
			}
			assert.Equalf(t, tt.want, s.AverageCredit(), "AverageCredit()")
		})
	}
}

func TestService_SummaryTransaction(t *testing.T) {
	type fields struct {
		idUser       int
		transactions internal.Transactions
		debit        internal.Transactions
		credit       internal.Transactions
		email        internal.EmailService
		userRep      internal.RepositoryUser
		transRep     internal.RepositoryTransaction
	}
	type args struct {
		ctx context.Context
	}

	transactions := helperTransactions(t)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Success getting the information and sent email",
			fields: fields{
				idUser:       1,
				transactions: transactions,
				email:        nil,
				userRep:      nil,
				transRep:     nil,
			},
			args:    args{},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				idUser:       tt.fields.idUser,
				transactions: tt.fields.transactions,
				debit:        tt.fields.debit,
				credit:       tt.fields.credit,
				email:        tt.fields.email,
				userRep:      tt.fields.userRep,
				transRep:     tt.fields.transRep,
			}
			tt.wantErr(t, s.SummaryTransaction(tt.args.ctx), fmt.Sprintf("SummaryTransaction(%v)", tt.args.ctx))
		})
	}
}

func helperTransactions(t *testing.T) internal.Transactions {
	transactionsArr := make([]internal.Transaction, 0)
	transactionsArr = append(transactionsArr, internal.Transaction{
		ID:          1,
		Date:        time.Now(),
		Transaction: 90,
	})
	transactionsArr = append(transactionsArr, internal.Transaction{
		ID:          2,
		Date:        time.Now(),
		Transaction: 10,
	})
	transactionsArr = append(transactionsArr, internal.Transaction{
		ID:          3,
		Date:        time.Now(),
		Transaction: -11,
	})
	return transactionsArr
}
