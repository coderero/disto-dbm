package types

// The Response struct is used to as a response body.
type Response struct {
	Status Status     `json:"status"`
	Data   any        `json:"data,omitempty"`
	Errors []APIError `json:"errors,omitempty"`
}

// The Status struct is used to as a response body.
type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type APIError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
