package appdef

type (
	Env struct {
		Dev        []EnvValue `json:"dev"`
		Staging    []EnvValue `json:"staging"`
		Production []EnvValue `json:"production"`
	}
	EnvValue struct {
		Key   string `json:"key"`
		Type  string `json:"type"`
		From  string `json:"from,omitempty"`
		Value string `json:"value,omitempty"`
	}
)
