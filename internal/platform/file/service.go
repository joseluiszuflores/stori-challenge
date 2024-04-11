package file

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	mooc "github.com/joseluiszuflores/stori-challenge/internal"

	"github.com/gocarina/gocsv"
)

type transaction struct {
	//nolint: revive,stylecheck
	Id          int
	Date        dateCSV
	Transaction float64
}

type dateCSV string

type transactions []transaction

type Service struct {
	pathCSV string
}

func NewService(pathCSV string) *Service {
	return &Service{pathCSV: pathCSV}
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
		newTimeParse, err := time.Parse(time.DateOnly, val.Date.DateWithYear())
		if err != nil {
			glog.Errorf("Error trying to add the date to [%v] [%s]", val, err)

			continue
		}
		mctrns = append(mctrns, mooc.Transaction{
			ID:          val.Id,
			Date:        newTimeParse,
			Transaction: val.Transaction,
		})
	}

	return mctrns
}

const separator = "/"

func (d dateCSV) Day() int {
	content := strings.Split(string(d), separator)
	if len(content) < lengthOfTimeIntoStruct {
		return 1
	}
	dayInt, err := strconv.Atoi(content[1])
	if err != nil {
		return 1
	}

	return dayInt
}

const lengthOfTimeIntoStruct = 2

func (d dateCSV) Month() int {
	content := strings.Split(string(d), separator)
	if len(content) < lengthOfTimeIntoStruct {
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
	month := fmt.Sprintf("%d", d.Month())
	if d.Month() < 10 {
		month = fmt.Sprintf("0%d", d.Month())
	}
	day := fmt.Sprintf("%d", d.Day())
	if d.Day() < 10 {
		day = fmt.Sprintf("0%d", d.Day())
	}

	return fmt.Sprintf("%d-%s-%s", d.Year(), month, day)
}
