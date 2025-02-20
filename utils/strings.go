package utils

import (
	"net/url"
	"regexp"
	"strings"

	"math/rand"
)

func Normalize(str string) string {
	str = strings.ToLower(str)

	charA := "àáãảạăằắẳẵặâầấẩẫậä"
	charE := "èéẻẽẹêềếểễệë"
	charU := "ùúủũụưừứửữựüû"
	charO := "òóỏõọôồốổỗộơờớởỡợö"
	charI := "ìíỉĩịïî"
	charY := "ýỳỹỵỷ"
	charD := "đ"

	for _, a := range charA {
		str = strings.ReplaceAll(str, string(a), "a")
	}

	for _, e := range charE {
		str = strings.ReplaceAll(str, string(e), "e")
	}

	for _, u := range charU {
		str = strings.ReplaceAll(str, string(u), "u")
	}

	for _, o := range charO {
		str = strings.ReplaceAll(str, string(o), "o")
	}

	for _, i := range charI {
		str = strings.ReplaceAll(str, string(i), "i")
	}

	for _, y := range charY {
		str = strings.ReplaceAll(str, string(y), "y")
	}

	for _, d := range charD {
		str = strings.ReplaceAll(str, string(d), "d")
	}

	return strings.ToUpper(str)
}

func TokenizeAndRejoinString(s string) string {
	tokens := strings.Fields(s) // Fields automatically splits on whitespace and removes extra spaces
	return strings.Join(tokens, " ")
}

func RemoveNewLines(input string) string {
	input = TokenizeAndRejoinString(input)
	// Replace literal backslash followed by 'n'
	result := strings.ReplaceAll(input, `\\n`, " ")
	// Replace double newline
	result = strings.ReplaceAll(result, "\n\n", " ")
	// Replace newline
	result = strings.ReplaceAll(result, "\n", " ")
	return result
}

func IsValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|<>?/"

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func AddUrlParams(rawURL, key, value string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	query := parsedURL.Query()
	query.Set(key, value)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}
