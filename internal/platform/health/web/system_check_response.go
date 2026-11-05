package web

// SystemCheckResponse defines a structured JSON format for status checking
type SystemCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}
