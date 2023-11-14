package types

// The SignUp struct is used to bind the request body form to the struct.
type Register struct {
	Username  string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,alpha"`
	LastName  string `json:"last_name" validate:"required,alpha"`
	Age       int    `json:"age" validate:"required,gt=0,lt=100"`
}

// The Login struct is used to bind the request body form to the struct.
type Login struct {
	Username string `json:"username" validate:"omitempty,alphanum"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required,min=8"`
}
