package utils_test

import (
	"testing"

	"github.com/ngtrvu/zen-go/utils"
	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "Nguyễn Văn Hải",
			expected: "NGUYEN VAN HAI",
		},
		{
			input:    "Nguyễn Văn nghiễm",
			expected: "NGUYEN VAN NGHIEM",
		},
		{
			input:    "Nguyễn Văn Khánh",
			expected: "NGUYEN VAN KHANH",
		},
		{
			input:    "Nguyễn Văn Lũ",
			expected: "NGUYEN VAN LU",
		},
		{
			input:    "Nguyễn Văn Lữ",
			expected: "NGUYEN VAN LU",
		},
		{
			input:    "Nguyễn Văn Tụ",
			expected: "NGUYEN VAN TU",
		},
		{
			input:    "Nguyễn Văn Tuỳ",
			expected: "NGUYEN VAN TUY",
		},
		{
			input:    "Huỳnh Bách Mỹ Duyên",
			expected: "HUYNH BACH MY DUYEN",
		},
		{
			input:    "Huỳnh Đạt",
			expected: "HUYNH DAT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := utils.Normalize(tt.input)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestAddUrlParams(t *testing.T) {
	tests := []struct {
		input    string
		key      string
		value    string
		expected string
	}{
		{
			input:    "https://example.com",
			key:      "view",
			value:    "mobile",
			expected: "https://example.com?view=mobile",
		},
		{
			input:    "https://example.com?x1=y1",
			key:      "view",
			value:    "mobile",
			expected: "https://example.com?view=mobile&x1=y1",
		},
		{
			input:    "https://example.com?x1=y1&x2=y2&view",
			key:      "view",
			value:    "mobile",
			expected: "https://example.com?view=mobile&x1=y1&x2=y2",
		},
		{
			input:    "https://example.com?x1=y1&x2=y2&view#test",
			key:      "view",
			value:    "mobile",
			expected: "https://example.com?view=mobile&x1=y1&x2=y2#test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := utils.AddUrlParams(tt.input, tt.key, tt.value)
			require.Equal(t, tt.expected, got)
		})
	}
}
