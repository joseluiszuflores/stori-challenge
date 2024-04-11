package file

import (
	"encoding/csv"
	"github.com/joseluiszuflores/stori-challenge/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestService_GetDataFile(t *testing.T) {
	type fields struct {
		pathCSV string
	}
	trans := helperFileCSV(t)
	tests := []struct {
		name    string
		fields  fields
		want    internal.Transactions
		wantErr bool
	}{
		{
			name:    "Success Open file tmp and reading the data",
			fields:  fields{pathCSV: filePathTest},
			want:    toMoocTransactions(trans),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				pathCSV: tt.fields.pathCSV,
			}
			got, err := s.GetDataFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

const filePathTest = "./file.csv"

func helperFileCSV(t *testing.T) transactions {
	t.Helper()
	trans := make(transactions, 0)
	trans = append(trans, transaction{
		Id:          1,
		Date:        "7/28",
		Transaction: +60.5,
	})
	trans = append(trans, transaction{
		Id:          2,
		Date:        "7/28",
		Transaction: +50.5,
	})

	trans = append(trans, transaction{
		Id:          2,
		Date:        "8/2",
		Transaction: +50.5,
	})
	file2, err := os.Create(filePathTest)
	if err != nil {
		t.Errorf("error creating the file [%s]", err)

		return nil
	}
	defer file2.Close()
	writer := csv.NewWriter(file2)
	defer writer.Flush()
	headers := []string{"Id", "Date", "Transaction"}
	assert.NoError(t, writer.Write(headers))
	assert.NoError(t, writer.Write([]string{"1", "7/28", "+61.5"}))
	assert.NoError(t, writer.Write([]string{"2", "7/28", "+50.5"}))
	return trans
}

func Test_dateCSV_DateWithYear(t *testing.T) {
	tests := []struct {
		name string
		d    dateCSV
		want string
	}{
		{
			name: "Success getting date",
			d:    dateCSV("7/15"),
			want: "2024/15/7",
		},
		{
			name: "Error getting date when the date is empty",
			d:    dateCSV(""),
			want: "2024/1/1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DateWithYear(); got != tt.want {
				t.Errorf("DateWithYear() = %v, want %v", got, tt.want)
			}
		})
	}
}
