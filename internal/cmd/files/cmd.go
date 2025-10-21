// Package files provides commands for generating and managing project configuration files
// such as package.json, turbo.json, and code style configurations.
package files

import (
	"encoding/json"
)

// identMarshaller is the function used to marshal JSON with indentation.
var identMarshaller = json.MarshalIndent
