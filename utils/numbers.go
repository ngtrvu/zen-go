package utils

import (
	"math"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func FormatCurrency(amount float64) string {
	p := message.NewPrinter(language.English)
	amountInt := int64(amount)
	return p.Sprintf("%d", amountInt)
}

func StringToFloat64(str string) *float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}
	return &val
}