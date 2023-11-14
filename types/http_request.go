package types

// RefreshToken is used to bind the request body form to the struct.
type RefreshToken struct {
	AcessToken   string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
