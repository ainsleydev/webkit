package payload

import (
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTime(t *testing.T) {
	t.Parallel()

	now := time.Now()
	got := NewTime(now)

	assert.Equal(t, now, got.Time())
	assert.Equal(t, now.Format(time.RFC3339), got.String())
}

func TestTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	input := Time{time: time.Date(2025, 1, 29, 12, 0, 0, 0, time.UTC)}
	want := `"2025-01-29T12:00:00Z"`

	got, err := json.Marshal(input)
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func TestTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		want    Time
		wantErr bool
	}{
		"Valid RFC3339 Time": {
			input:   `"2025-01-29T12:00:00Z"`,
			want:    Time{time: time.Date(2025, 1, 29, 12, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		"Invalid Time Format": {
			input:   `"invalid-time"`,
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var got Time
			err := json.Unmarshal([]byte(test.input), &got)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want.Time(), got.Time())
		})
	}
}

func TestTime_Time(t *testing.T) {
	t.Parallel()
	in := time.Date(2025, 1, 29, 12, 0, 0, 0, time.UTC)
	input := Time{time: in}
	assert.Equal(t, in, input.Time())
}

func TestTime_String(t *testing.T) {
	t.Parallel()
	in := "2025-01-29T12:00:00Z"
	input := Time{val: in}
	assert.Equal(t, in, input.String())
}

func TestTime_IsZero(t *testing.T) {
	type Test struct {
		CheckedAt *Time `json:"checkedAt,omitempty"`
		CreatedAt *Time `json:"createdAt,omitempty"`
		UpdatedAt *Time `json:"updatedAt,omitempty"`
	}
	test := Test{}
	data, err := json.Marshal(test)
	require.NoError(t, err)
	assert.Equal(t, `{}`, string(data))
}
