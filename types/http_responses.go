package types

// The ResponseSkeleton type is used to represent a skeleton of a response in Go.
type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}
