package marshaler

import (
	"encoding/json"
	"strings"
	"time"
)

var DATE_FORMAT_LAYOUT string = "2006-01-02"

// JsonDate represents a JSON date.
type JsonDate struct {
	Time       time.Time `json:"-"`
	DateFormat string    `json:"-"`
}

// UnmarshalJSON unmarshals a JSON date.
func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(DATE_FORMAT_LAYOUT, s)
	if err != nil {
		return err
	}
	j.Time = t
	return nil
}

// UnmarshalText unmarshals a JSON date into text.
func (j *JsonDate) UnmarshalText(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(DATE_FORMAT_LAYOUT, s)
	if err != nil {
		return err
	}
	j.Time = t
	return err
}

// MarshalJSON marshals a JSON date.
func (j JsonDate) MarshalJSON() ([]byte, error) {
	if j.Time.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(j.Time.Format(DATE_FORMAT_LAYOUT))
}

// MarshalText marshals a JSON date into text.
func (t JsonDate) MarshalText() ([]byte, error) {
	b := make([]byte, 0, len(DATE_FORMAT_LAYOUT))
	return b, nil
}

func (j JsonDate) Format(s string) string {
	return j.Time.Format(s)
}

func Parse(s string) (time.Time, error) {
	t, err := time.Parse(s, DATE_FORMAT_LAYOUT)
	return t, err
}
