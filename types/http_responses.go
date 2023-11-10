package types

// The ResponseSkeleton type is used to represent a skeleton of a response in Go.
type Response struct {
	Status Status         `json:"status"`
	Data   map[string]any `json:"data,omitempty"`
}

type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
