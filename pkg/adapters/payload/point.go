package payload

// Point represents a geographical coordinate with
// longitude and latitude, stored as a slice of two float64 values where:
//
// - The first element (index 0) is the longitude.
// - The second element (index 1) is the latitude.
//
// See: https://payloadcms.com/docs/beta/fields/point
type Point []float64

// Latitude returns the latitude value from the Point.
// If the Point is invalid (does not contain two elements), it returns 0.
func (p Point) Latitude() float64 {
	if len(p) > 1 {
		return p[1]
	}
	return 0
}

// Longitude returns the longitude value from the Point.
// If the Point is invalid (does not contain two elements), it returns 0.
func (p Point) Longitude() float64 {
	if len(p) > 0 {
		return p[0]
	}
	return 0
}
