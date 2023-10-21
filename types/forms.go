package types

type Register struct {
	Email     string `json:"email" form:"email" binding:"required,email"`
	Password  string `json:"password" form:"password" binding:"required,min=8,max=32"`
	Username  string `json:"username" form:"username" binding:"required,min=4,max=32"`
	FirstName string `json:"first_name" form:"firstname" binding:"required,min=4,max=32"`
	LastName  string `json:"last_name" form:"lastname" binding:"required,min=4,max=32"`
	Age       int    `json:"age" form:"age" binding:"required,min=1,max=3"`
}

type Login struct {
	Email    string `json:"email" form:"email" binding:"email"`
	Username string `json:"username" form:"username" binding:"min=4,max=32"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=32"`
}
