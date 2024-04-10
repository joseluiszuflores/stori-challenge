package file

import (
	"fmt"
	"github.com/golang/glog"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type transaction struct {
	Id          int
	Date        dateCSV
	Transaction float64
}

type dateCSV string

type transactions []transaction

type Service struct {
	pathCSV string
}

func (s *Service) GetDataFile() (mooc.Transactions, error) {
	file, err := os.Open(s.pathCSV)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tns := make(transactions, 0)
	err = gocsv.UnmarshalFile(file, &tns)
	if err != nil {
		return nil, err
	}
	return toMoocTransactions(tns), nil
}

func toMoocTransactions(ts transactions) mooc.Transactions {
	mctrns := make(mooc.Transactions, 0)
	for _, val := range ts {
		tm, err := time.Parse(time.DateOnly, val.Date.DateWithYear())
		if err != nil {
			glog.Errorf("Error trying to add the date to [%v] [%s]", val, err)
			continue
		}
		mctrns = append(mctrns, mooc.Transaction{
			ID:          val.Id,
			Date:        tm,
			Transaction: val.Transaction,
		})
	}

	return mctrns
}

const separator = "/"

func (d dateCSV) Day() int {
	content := strings.Split(string(d), separator)
	if len(content) < 2 {
		return 1
	}
	dayInt, err := strconv.Atoi(content[1])
	if err != nil {
		return 1
	}
	return dayInt
}

func (d dateCSV) Month() int {
	content := strings.Split(string(d), separator)
	if len(content) < 2 {
		return 1
	}
	mothInt, err := strconv.Atoi(content[0])
	if err != nil {
		return 1
	}
	return mothInt
}

func (d dateCSV) Year() int {
	return time.Now().Year()
}

func (d dateCSV) DateWithYear() string {
	return fmt.Sprintf("%d-0%d-%d", d.Year(), d.Month(), d.Day())
}
