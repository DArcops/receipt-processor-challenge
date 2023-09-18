package util

import "time"

const (
	dateLayout     = "2006-01-02"
	hourTimeLayout = "15:04"
)

// IsTimeBetween checks if a time is between two other times.
func IsTimeBetween(timeToCheck string, startTimeHour, endTimeHour int) (bool, error) {
	date, err := time.Parse(hourTimeLayout, timeToCheck)
	if err != nil {
		return false, err
	}

	st := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		startTimeHour,
		0,
		0,
		0,
		date.Location(),
	)

	et := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		endTimeHour,
		0,
		0,
		0,
		date.Location(),
	)

	return date.Equal(st) || (date.After(st) && date.Before(et)) || date.Equal(et), nil
}

// IsDayOdd checks if a day is odd.
func IsDayOdd(date string) (bool, error) {
	formatedDate, err := time.Parse(dateLayout, date)
	if err != nil {
		return false, err
	}

	day := formatedDate.Day()

	return day%2 != 0, nil
}
