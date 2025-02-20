package marshaler_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ngtrvu/zen-go/marshaler"
	"github.com/stretchr/testify/assert"
)

type TestSerializer struct {
	DateOfBirth marshaler.JsonDate `json:"date_of_birth,omitempty"`
}

func TestJsonDate(t *testing.T) {
	var serializerData TestSerializer
	payload := strings.NewReader(`{"date_of_birth": "2000-09-03"}`)
	json.NewDecoder(payload).Decode(&serializerData)

	data, _ := json.Marshal(serializerData)
	assert.Equal(t, 3, serializerData.DateOfBirth.Time.Day())
	assert.Equal(t, time.Month(9), serializerData.DateOfBirth.Time.Month())
	assert.Equal(t, 2000, serializerData.DateOfBirth.Time.Year())

	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	assert.Equal(t, "2000-09-03", result["date_of_birth"])

	testTime, _ := time.Parse("2006-01-02", "2000-09-03")
	s := &TestSerializer{DateOfBirth: marshaler.JsonDate{Time: testTime}}
	b, _ := json.Marshal(s)
	assert.Equal(t, "{\"date_of_birth\":\"2000-09-03\"}", string(b))

	s2 := &TestSerializer{DateOfBirth: marshaler.JsonDate{Time: time.Time{}}}
	b2, _ := json.Marshal(s2)
	assert.Equal(t, "{\"date_of_birth\":\"\"}", string(b2))
}
