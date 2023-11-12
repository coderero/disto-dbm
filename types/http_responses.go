package types

// The Response struct is used to as a response body.
type Response struct {
	Status Status         `json:"status"`
	Data   map[string]any `json:"data,omitempty"`
}

// The Status struct is used to as a response body.
type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
