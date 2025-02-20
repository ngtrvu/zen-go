package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ngtrvu/zen-go/utils"
)

func TestMath_RoundFloat(t *testing.T) {
	x1 := utils.RoundFloat(1.23456789, 2)
	expected1 := 1.23
	assert.Equal(t, expected1, x1)

	x2 := utils.RoundFloat(1.23456789, 0)
	expected2 := 1.0
	assert.Equal(t, expected2, x2)

	x3 := utils.RoundFloat(1.23456789, 3)
	expected3 := 1.235
	assert.Equal(t, expected3, x3)

	x4 := utils.RoundFloat(1.23456789, 4)
	expected4 := 1.2346
	assert.Equal(t, expected4, x4)

	x5 := utils.RoundFloat(1.0000000000001, 2)
	expected5 := 1.0
	assert.Equal(t, expected5, x5)

	x6 := utils.RoundFloat(10000.0000000000001, 0)
	expected6 := 10000.0
	assert.Equal(t, expected6, x6)

	x7 := utils.RoundFloat(10000.499999999, 0)
	expected7 := 10000.0
	assert.Equal(t, expected7, x7)

	x8 := utils.RoundFloat(10000.599999999, 0)
	expected8 := 10001.0
	assert.Equal(t, expected8, x8)

	x9 := utils.RoundFloat(994961.7312, 0)
	expected9 := 994962.0
	assert.Equal(t, expected9, x9)

}

func Test_FormatCurrency(t *testing.T) {
	x1 := utils.FormatCurrency(1000000)
	expected1 := "1,000,000"
	assert.Equal(t, expected1, x1)

	x2 := utils.FormatCurrency(100000)
	expected2 := "100,000"
	assert.Equal(t, expected2, x2)

	x3 := utils.FormatCurrency(9999999)
	expected3 := "9,999,999"
	assert.Equal(t, expected3, x3)

	x4 := utils.FormatCurrency(500000000000)
	expected4 := "500,000,000,000"
	assert.Equal(t, expected4, x4)

	x5 := utils.FormatCurrency(500)
	expected5 := "500"
	assert.Equal(t, expected5, x5)
}
