package utils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ngtrvu/zen-go/utils"
)

func TestTime_GetNowDate(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")

	expected := time.Now().In(location).Format("2006/01/02")
	nowDate := utils.GetNowInputDate("")
	assert.Equal(t, expected, nowDate)

	expected = time.Now().In(location).Format("2006/01/02")
	nowDate = utils.GetNowInputDate("2006/01/02")
	assert.Equal(t, expected, nowDate)

	expected = time.Now().In(location).Format("2006-01-02")
	nowDate = utils.GetNowInputDate("2006-01-02")
	assert.Equal(t, expected, nowDate)

}

func TestTime_GetPreviousDate(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	nDays := 3
	expected := time.Now().In(location).AddDate(0, 0, -nDays).Format("2006/01/02")
	previousDate := utils.GetPreviousDate(nDays, "2006/01/02")
	assert.Equal(t, expected, previousDate)

	expected = time.Now().In(location).Format("2006/01/02")
	previousDate = utils.GetPreviousDate(0, "2006/01/02")
	assert.Equal(t, expected, previousDate)

	expected = time.Now().In(location).Format("2006-01-02")
	previousDate = utils.GetPreviousDate(0, "2006-01-02")
	assert.Equal(t, expected, previousDate)
}

func TestTime_GetDateInPreviousYears(t *testing.T) {
	nYears := 5
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	expected := time.Now().In(location).AddDate(-nYears, 0, 0).Format("2006/01/02")
	previousYearDate := utils.GetPreviousYearDateFromNow(nYears, "2006/01/02")
	assert.Equal(t, expected, previousYearDate)

	expected = time.Now().In(location).AddDate(-nYears, 0, 0).Format("2006/01/02")
	previousYearDate = utils.GetPreviousYearDateFromNow(nYears, "2006/01/02")
	assert.Equal(t, expected, previousYearDate)

	expected = time.Now().In(location).AddDate(-nYears, 0, 0).Format("2006-01-02")
	previousYearDate = utils.GetPreviousYearDateFromNow(nYears, "2006-01-02")
	assert.Equal(t, expected, previousYearDate)
}

func TestTime_GetNextDate(t *testing.T) {
	nDays := 3
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	expected := time.Now().In(location).AddDate(0, 0, nDays).Format("2006/01/02")
	nextDate := utils.GetNextDate(nDays, "2006/01/02")
	assert.Equal(t, expected, nextDate)

	expected = time.Now().In(location).Format("2006/01/02")
	nextDate = utils.GetNextDate(0, "2006/01/02")
	assert.Equal(t, expected, nextDate)

	expected = time.Now().In(location).Format("2006-01-02")
	nextDate = utils.GetNextDate(0, "2006-01-02")
	assert.Equal(t, expected, nextDate)
}

func TestTime_GetCurrentDatePath(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	expected := time.Now().In(location).Format("2006/01/02")
	currentDate := utils.GetCurrentDatePath()

	assert.Equal(t, expected, currentDate)
}

func TestTime_BuildCurrentDatePath(t *testing.T) {
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	expected := "your/path/here/" + time.Now().In(location).Format("2006/01/02")
	currentDate := utils.BuildCurrentDatePath("your/path/here")

	assert.Equal(t, expected, currentDate)
}

func TestTime_ComputeDiffYears(t *testing.T) {
	nYears := 5
	location, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	expected := time.Now().In(location).AddDate(-nYears, 0, 0).Format("2006/01/02")
	previousYearDate := utils.GetPreviousYearDateFromNow(nYears, "2006/01/02")
	assert.Equal(t, expected, previousYearDate)
	actualNYears := utils.ComputeDiffYears(time.Now().In(location), time.Now().In(location).AddDate(-nYears, 0, 0))
	assert.Equal(t, nYears, actualNYears)

	actualNYears = utils.ComputeDiffYears(time.Now().In(location), time.Now().In(location).AddDate(-nYears, +1, 0))
	assert.Equal(t, nYears-1, actualNYears)

	actualNYears = utils.ComputeDiffYears(time.Now().In(location), time.Now().In(location).AddDate(-nYears, 0, 0))
	assert.Equal(t, nYears, actualNYears)

}
