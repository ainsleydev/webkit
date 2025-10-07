package appdef

import (
	"encoding/json"
	"fmt"

	"github.com/swaggest/jsonschema-go"
)

// CommandSpec defines an action for an App, this can run in
// CI or locally.
type CommandSpec struct {
	Cmd      string `json:"command,omitempty"`
	SkipCI   bool   `json:"skip_ci,omitempty"`
	Timeout  string `json:"timeout,omitempty"`
	Disabled bool   `json:"-"` // Set during unmarshal
}

// Ensure CommandSpec implements jsonschema.OneOfExposer
var _ jsonschema.OneOfExposer = &CommandSpec{}

// JSONSchemaOneOf returns the polymorphic schema options.
func (*CommandSpec) JSONSchemaOneOf() []interface{} {
	return []interface{}{
		true,          // boolean option
		"",            // string option
		CommandSpec{}, // object option
	}
}

// Command defines the type of action that will be actioned.
type Command string

// Command constants.
const (
	CommandLint   Command = "lint"
	CommandTest   Command = "test"
	CommandFormat Command = "format"
)

// commands defines all the commands available.
var commands = []Command{
	CommandLint,
	CommandTest,
	CommandFormat,
}

// defaultCommands defines the default actions for each
// application type. If not overridden, these commands
// will be used.
var defaultCommands = map[AppType]map[Command]string{
	AppTypePayload: {
		CommandLint:   "pnpm lint",
		CommandTest:   "pnpm test",
		CommandFormat: "pnpm format",
	},
	AppTypeSvelteKit: {
		CommandLint:   "pnpm lint",
		CommandTest:   "pnpm test",
		CommandFormat: "pnpm format",
	},
	AppTypeGoLang: {
		CommandLint:   "golangci-lint run",
		CommandTest:   "go test ./...",
		CommandFormat: "gofmt -w .",
	},
}

// UnmarshalJSON implements json.Unmarshaler to
func (c *CommandSpec) UnmarshalJSON(data []byte) error {
	// Try bool (false = disabled)
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		c.Disabled = !b
		return nil
	}

	// Try string (override)
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		c.Cmd = str
		return nil
	}

	// Try object (full control)
	type Alias CommandSpec
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("invalid command format: expected bool, string, or object")
	}

	*c = CommandSpec(aux)
	return nil
}

// String implements fmt.Stringer on Command.
func (c Command) String() string {
	return string(c)
}
