package types

// APIInfo represents API metadata and information
type APIInfo struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Version     string      `json:"version"`
	Status      string      `json:"status"`
	Environment string      `json:"environment"`
	Endpoints   *Endpoints  `json:"endpoints,omitempty"`
	Links       *APILinks   `json:"links,omitempty"`
}

// Endpoints contains available API endpoints
type Endpoints struct {
	Health     string `json:"health"`
	HealthDB   string `json:"health_db"`
	API        string `json:"api"`
	Documentation string `json:"documentation,omitempty"`
}

// APILinks contains useful links
type APILinks struct {
	Self        string `json:"self"`
	Health      string `json:"health"`
	HealthDB    string `json:"health_db"`
	Documentation string `json:"documentation,omitempty"`
}
