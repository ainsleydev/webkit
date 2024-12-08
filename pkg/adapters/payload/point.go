package payload

import (
	"encoding/json"
	"fmt"
)

// Point represents a geographical coordinate.
// While it's defined as a struct, it marshals to/from JSON as a [longitude, latitude] array.
type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// MarshalJSON implements the json.Marshaler interface.
// Converts the Point to a [longitude, latitude] array.
//
//goland:noinspection GoMixedReceiverTypes
func (p Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToSlice())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Expects a [longitude, latitude] array.
//
//goland:noinspection GoMixedReceiverTypes
func (p *Point) UnmarshalJSON(data []byte) error {
	var coords []float64
	if err := json.Unmarshal(data, &coords); err != nil {
		return err
	}
	if len(coords) != 2 {
		return fmt.Errorf("point array must contain exactly 2 elements [longitude, latitude]")
	}
	p.Longitude = coords[0]
	p.Latitude = coords[1]
	return nil
}

// ToSlice converts the Point to a []float64 slice in [longitude, latitude] order.
//
//goland:noinspection GoMixedReceiverTypes
func (p Point) ToSlice() []float64 {
	return []float64{p.Longitude, p.Latitude}
}
