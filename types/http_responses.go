package types

// The SuccessResponse type represents a successful response with a ResponseSkeleton.
// @property {ResponseSkeleton} Success - The `Success` property is of type `ResponseSkeleton` and is
// used to represent a successful response in a JSON format.
type SuccessResponse struct {
	Success ResponseSkeleton `json:"success"`
}

// The ErrorResponse type is used to represent an error response in Go.
// @property {ResponseSkeleton} Error - The `Error` property is of type `ResponseSkeleton`.
type ErrorResponse struct {
	Error ResponseSkeleton `json:"error"`
}

// The ResponseSkeleton type is used to represent a skeleton of a response in Go.
type ResponseSkeleton struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}
