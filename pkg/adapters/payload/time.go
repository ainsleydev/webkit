package payload

import (
	"time"

	"github.com/goccy/go-json"
)

// Time is a wrapper around time.Time to handle RFC3339 formatting
// from Payload, as it returns times as strings.
type Time struct {
	val  string
	time time.Time
}

// NewTime creates a new Payload Time type from the given input.
func NewTime(t time.Time) *Time {
	return &Time{
		val:  t.Format(time.RFC3339),
		time: t,
	}
}

// MarshalJSON converts Time to a JSON string in RFC3339 format.
//
//goland:noinspection GoMixedReceiverTypes
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.time.Format(time.RFC3339))
}

// UnmarshalJSON parses a JSON string into a Time object
//
//goland:noinspection GoMixedReceiverTypes
func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}
	t.time = parsedTime
	t.val = s
	return nil
}

// Time returns the underlying time.Time value.
//
//goland:noinspection GoMixedReceiverTypes
func (t Time) Time() time.Time {
	return t.time
}

// String implements fmt.Stringer to return the time
// value as a RFC3339 formatted string.
//
//goland:noinspection GoMixedReceiverTypes
func (t Time) String() string {
	return t.val
}

// IsZero returns true if the Time struct should be treated as empty.
//
//goland:noinspection GoMixedReceiverTypes
func (t Time) IsZero() bool {
	return t.val == "" || t.time.IsZero()
}
