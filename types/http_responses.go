package types

// The ResponseSkeleton type is used to represent a skeleton of a response in Go.
type Response struct {
	Status     bool           `json:"status"`
	StatusCode int            `json:"status_code,omitempty"`
	Message    string         `json:"message"`
	Data       map[string]any `json:"data,omitempty"`
}
