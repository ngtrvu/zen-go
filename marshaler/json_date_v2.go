package marshaler

import (
	"fmt"
	"strings"
	"time"
)

// JsonDate represents a JSON date.
// swagger:strfmt date
type JsonDateV2 time.Time

// UnmarshalJSON unmarshals a JSON date.
func (jd *JsonDateV2) UnmarshalJSON(input []byte) error {
	// Trim the quotes and unmarshal as time.Time
	strInput := string(input)
	strInput = strInput[1 : len(strInput)-1]

	newTime, err := time.Parse(DATE_FORMAT_LAYOUT, strInput)
	if err != nil {
		return err
	}

	*jd = JsonDateV2(newTime)
	return nil
}

// UnmarshalText unmarshals a JSON date into text.
func (jd *JsonDateV2) UnmarshalText(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(DATE_FORMAT_LAYOUT, s)
	if err != nil {
		return err
	}

	*jd = JsonDateV2(t)
	return nil
}

func (jd JsonDateV2) MarshalJSON() ([]byte, error) {
	// Format the time as a string and add quotes
	str := fmt.Sprintf("\"%s\"", time.Time(jd).Format(DATE_FORMAT_LAYOUT))

	// Convert the string to a byte slice
	return []byte(str), nil
}

func (jd JsonDateV2) MarshalText() ([]byte, error) {
	// Format the time as a string
	str := time.Time(jd).Format(DATE_FORMAT_LAYOUT)

	// Convert the string to a byte slice
	return []byte(str), nil
}
