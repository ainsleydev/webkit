package appdef

type (
	Resource struct {
		Name     string         `json:"name"`
		Type     string         `json:"type"`
		Provider string         `json:"provider"`
		Config   map[string]any `json:"config"` // Conforms to Terraform
		Outputs  []string       `json:"outputs"`
		Backup   any            // TODO
	}
)
