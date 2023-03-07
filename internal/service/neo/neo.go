package neo

import (
	"errors"
	"time"
)

var ErrNoResult = errors.New("neo not found")

type Neo struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func NewNeo(date time.Time, count int) Neo {
	return Neo{
		Date:  date,
		Count: count,
	}
}

func ParseDates(dates []string) ([]time.Time, error) {
	var parsedDates []time.Time
	for _, date := range dates {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			return nil, err
		}
		parsedDates = append(parsedDates, parsedDate)
	}
	if len(parsedDates) < 1 {
		return nil, errors.New("dates not found")
	}
	return parsedDates, nil
}

func GetDefaultDatesForScan() []time.Time {
	today := time.Now()
	dates := make([]time.Time, 0, 7)
	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, i)
		dates = append(dates, date)
	}
	return dates
}
