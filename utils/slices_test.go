package utils_test

import (
	"reflect"
	"testing"

	"github.com/ngtrvu/zen-go/utils"
)

func TestRemoveInt(t *testing.T) {
	tests := []struct {
		name string
		s    []int
		item int
		want []int
	}{
		{
			name: "remove item from middle",
			s:    []int{1, 2, 3, 4, 5},
			item: 3,
			want: []int{1, 2, 4, 5},
		},
		{
			name: "remove item from end",
			s:    []int{1, 2, 3, 4, 5},
			item: 5,
			want: []int{1, 2, 3, 4},
		},
		{
			name: "remove item not in slice",
			s:    []int{1, 2, 3, 4, 5},
			item: 6,
			want: []int{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Remove(tt.s, tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveString(t *testing.T) {
	tests := []struct {
		name string
		s    []string
		item string
		want []string
	}{
		{
			name: "remove item from middle",
			s:    []string{"a", "b", "c", "d", "e"},
			item: "c",
			want: []string{"a", "b", "d", "e"},
		},
		{
			name: "remove item from end",
			s:    []string{"a", "b", "c", "d", "e"},
			item: "e",
			want: []string{"a", "b", "c", "d"},
		},
		{
			name: "remove item not in slice",
			s:    []string{"a", "b", "c", "d", "e"},
			item: "f",
			want: []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Remove(tt.s, tt.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomCode(t *testing.T) {
	code1 := utils.GenerateRandomCode()
	if len(code1) != 6 {
		t.Errorf("GenerateRandomCode() = %v, want %v", code1, "6")
	}

	code2 := utils.GenerateRandomCode()

	if code1 == code2 {
		t.Errorf("generate random code is not unique. should be unique even if call at same time")
	}
}

type ExampleStruct struct {
	Name   string
	Status int
}

func TestFilterStruct(t *testing.T) {
	tests := []struct {
		name        string
		s           []ExampleStruct
		filterField string
		filterValue any
		want        []ExampleStruct
	}{
		{
			name: "filter status = 1",
			s: []ExampleStruct{
				{Name: "a", Status: 1},
				{Name: "b", Status: 2},
				{Name: "c", Status: 1},
			},
			filterField: "Status",
			filterValue: 1,
			want: []ExampleStruct{
				{Name: "a", Status: 1},
				{Name: "c", Status: 1},
			},
		},
		{
			name: "filter name = a",
			s: []ExampleStruct{
				{Name: "a", Status: 1},
				{Name: "b", Status: 2},
				{Name: "a", Status: 3},
			},
			filterField: "Name",
			filterValue: "a",
			want: []ExampleStruct{
				{Name: "a", Status: 1},
				{Name: "a", Status: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.Filter(tt.s, tt.filterField, tt.filterValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValuesString(t *testing.T) {
	tests := []struct {
		name         string
		s            []ExampleStruct
		filterField  string
		exampleValue string
		want         []string
	}{
		{
			name: "get values field name",
			s: []ExampleStruct{
				{Name: "a", Status: 1},
				{Name: "b", Status: 2},
				{Name: "c", Status: 3},
			},
			filterField:  "Name",
			exampleValue: "string",
			want:         []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GetValues(tt.s, tt.filterField, tt.exampleValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetValuesInt(t *testing.T) {
	tests := []struct {
		name         string
		s            []ExampleStruct
		filterField  string
		exampleValue int
		want         []int
	}{
		{
			name: "get values field Status",
			s: []ExampleStruct{
				{Name: "a", Status: 1},
				{Name: "b", Status: 2},
				{Name: "c", Status: 3},
			},
			filterField:  "Status",
			exampleValue: 1,
			want:         []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.GetValues(tt.s, tt.filterField, tt.exampleValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
