package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// Package-level constants (accessible throughout the package)
const (
	UtcTimezone     = "UTC"
	VnTimezone      = "Asia/Ho_Chi_Minh"
	FormatInputDate = "2006/01/02"
)

func GetNowInputDate(dateFormat string) string {
	location, _ := time.LoadLocation(VnTimezone)
	if dateFormat == "" {
		dateFormat = FormatInputDate
	}
	inputDate := time.Now().In(location).Format(dateFormat)

	return inputDate
}

func GetNow() time.Time {
	location, _ := time.LoadLocation(VnTimezone)
	now := time.Now().In(location)
	return now
}

func GetBeginOfDaytime() time.Time {
	t := GetNow()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetPreviousDate(nDay int, dateFormat string) string {
	location, _ := time.LoadLocation(VnTimezone)
	if dateFormat == "" {
		dateFormat = FormatInputDate
	}
	// To subtract days, use AddDate(0, 0, -days).
	nPreviousDateTime := time.Now().In(location).AddDate(0, 0, -nDay)
	previousDate := nPreviousDateTime.Format(dateFormat)

	return previousDate
}

func IsSameDay(date1 time.Time, date2 time.Time) bool {
	return date1.Year() == date2.Year() && date1.YearDay() == date2.YearDay()
}

func GetNextDate(nDay int, dateFormat string) string {
	location, _ := time.LoadLocation(VnTimezone)
	if dateFormat == "" {
		dateFormat = FormatInputDate
	}
	nextDateTime := time.Now().In(location).AddDate(0, 0, nDay)
	nextDate := nextDateTime.Format(dateFormat)

	return nextDate
}

func GetPreviousYearDateFromNow(nYear int, dateFormat string) string {
	location, _ := time.LoadLocation(VnTimezone)
	if dateFormat == "" {
		dateFormat = FormatInputDate
	}
	previousYearDateTime := time.Now().In(location).AddDate(-nYear, 0, 0)
	previousYearDate := previousYearDateTime.Format(dateFormat)

	return previousYearDate
}

func GetDatePath(date time.Time) string {
	dateFormat := "2006/01/02"
	return date.Format(dateFormat)
}

func GetCurrentDatePath() string {
	location, _ := time.LoadLocation(VnTimezone)
	return GetDatePath(time.Now().In(location))
}

func BuildCurrentDatePath(path string) string {
	return strings.Join([]string{path, GetCurrentDatePath()}, "/")
}

func ReformatDateString(dateStr string, oldLayout string, newLayout string, timezone string) (string, error) {
	t, err := time.Parse(oldLayout, dateStr)
	if err != nil {
		return "", fmt.Errorf("time.Parse: %v", err)
	}
	if timezone == "" {
		timezone = VnTimezone
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", fmt.Errorf("time.LoadLocation: %v", err)
	}
	t = t.In(loc)

	return t.Format(newLayout), nil
}

func DateStringToTimestamp(dateString string, dateFormat string) (float64, error) {
	if dateString == "" {
		dateString = GetNowInputDate("2006-01-02")
	}
	if dateFormat == "" {
		dateFormat = "2006-01-02"
	}
	location, _ := time.LoadLocation(VnTimezone)
	dateObject, err := time.ParseInLocation(dateFormat, dateString, location)
	if err != nil {
		return math.NaN(), err
	}
	timestamp := dateObject.Unix()

	return float64(timestamp), nil
}

func ComputeDiffYears(toDate time.Time, fromDate time.Time) int {
	years := toDate.Year() - fromDate.Year()
	months := int(toDate.Month()) - int(fromDate.Month())

	// If the month of toDate is before the month of fromDate, subtract one year
	if months < 0 || (months == 0 && toDate.Day() < fromDate.Day()) {
		years--
	}

	return years
}

func ComputeDiffDays(toDate time.Time, fromDate time.Time) int {
	duration := toDate.Sub(fromDate)
	days := int(duration.Hours() / 24)
	return days
}
