package appdef

import (
	"encoding/json"
	"fmt"

	"github.com/swaggest/jsonschema-go"
)

// CommandSpec defines an action for an App, this can run in
// CI or locally.
type CommandSpec struct {
	Name     string `json:"-"`
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
	CommandBuild  Command = "build"
)

// Commands defines all the Commands available that should be
// run in order.
var Commands = []Command{
	CommandFormat,
	CommandLint,
	CommandTest,
	CommandBuild,
}

// defaultCommands defines the default actions for each
// application type. If not overridden, these Commands
// will be used.
var defaultCommands = map[AppType]map[Command]string{
	AppTypePayload: {
		CommandFormat: "pnpm format",
		CommandLint:   "pnpm lint",
		CommandTest:   "pnpm test",
		CommandBuild:  "pnpm build",
	},
	AppTypeSvelteKit: {
		CommandFormat: "pnpm format",
		CommandLint:   "pnpm lint",
		CommandTest:   "pnpm test",

		CommandBuild: "pnpm build",
	},
	AppTypeGoLang: {
		CommandFormat: "gofmt -w .",
		CommandLint:   "golangci-lint run",
		CommandTest:   "go test ./...",
		CommandBuild:  "go build main.go",
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
