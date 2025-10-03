package app

type Definition struct {
	Project Project `json:"project"`
}

type Project struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Repository  string `json:"repository"`
}
