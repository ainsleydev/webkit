package appdef

// Config is a generic configuration map used across apps, resources, and monitors.
// It provides type-safe accessor methods for common configuration value types.
type Config map[string]any

// String safely retrieves a string value from the config.
// Returns the value and true if found, empty string and false otherwise.
func (c Config) String(key string) (string, bool) {
	if c == nil {
		return "", false
	}
	val, ok := c[key].(string)
	return val, ok
}

// Int safely retrieves an int value from the config.
// Returns the value and true if found, 0 and false otherwise.
func (c Config) Int(key string) (int, bool) {
	if c == nil {
		return 0, false
	}
	val, ok := c[key].(int)
	return val, ok
}

// Bool safely retrieves a bool value from the config.
// Returns the value and true if found, false and false otherwise.
func (c Config) Bool(key string) (bool, bool) {
	if c == nil {
		return false, false
	}
	val, ok := c[key].(bool)
	return val, ok
}
