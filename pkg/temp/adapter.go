package temp

// Adapter for on different platforms such as Payload & Static
type Adapter interface {
	Head() string

	Redirect()
}
