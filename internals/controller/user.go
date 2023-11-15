package controller

import (
	"net/http"
	"strings"

	"coderero.dev/projects/go/gin/hello/cache"
	"coderero.dev/projects/go/gin/hello/models"
	"coderero.dev/projects/go/gin/hello/pkg/security"
	"coderero.dev/projects/go/gin/hello/pkg/utils"
	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// The `PartialUpdate` struct is used to validate the request body for the `UpdateUser` function.
type PartialUpdate struct {
	Username    string `json:"username" validate:"omitempty,min=3,max=32,alphanum"`
	Email       string `json:"email" validate:"omitempty,email"`
	Password    string `json:"password" validate:"omitempty,min=8"`
	NewPassword string `json:"new_password" validate:"omitempty,min=8"`
	FirstName   string `json:"first_name" validate:"omitempty,alpha"`
	LastName    string `json:"last_name" validate:"omitempty,alpha"`
	Age         int    `json:"age" validate:"omitempty,gt=0,lt=100"`
}

// The `GetUser` function is a handler function for retrieving a user's information. It takes a
// `gin.Context` object as a parameter, which provides access to the request and response objects.
func (UserController) GetUser(c *gin.Context) {
	// The code `sub, _, shouldReturn := getSubFromToken(c)` is calling the `getSubFromToken` function and
	// assigning the returned values to the variables `sub`, `_`, and `shouldReturn`.
	sub, _, shouldReturn := getSubFromToken(c)
	if shouldReturn {
		return
	}

	// The code is declaring a variable `user` of type `*models.User` and then calling the `GetUserByEmail`
	// method on the `user` variable.
	var user *models.User

	err := user.GetUserByEmail(sub)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "user not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "ok",
		},
		Data: []map[string]interface{}{
			{
				"username":   user.Username,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"age":        user.Age,
			},
		}})
}

// The `func (*UserController) UpdateUser(c *gin.Context)` function is a method of the `UserController` struct which is
// used as a handler function for updating a user.
func (UserController) UpdateUser(c *gin.Context) {
	sub, _, shouldReturn := getSubFromToken(c)
	if shouldReturn {
		return
	}

	// The code is declaring a variable `registeredUser` of type `*models.User` and then assigning the
	// result of calling the `GetUserByEmail` method on the `registeredUser` variable. The `GetUserByEmail`
	// method is a function that retrieves a user from the database based on the provided email. The `sub`
	// variable is passed as an argument to the `GetUserByEmail` method to retrieve the user associated
	// with that email.
	var registeredUser *models.User
	registeredUser.GetUserByEmail(sub)

	// The code is declaring a variable `partialUpdate` of type `PartialUpdate` and then calling the
	// `DecodeJson` function from the `utils` package. The `DecodeJson` function is used to decode the JSON
	// request body from the `gin.Context` `c` and populate the `partialUpdate` variable with the decoded
	// values. If there is an error during the decoding process, the function will return and exit the
	// `UpdateUser` handler function.
	var partialUpdate PartialUpdate
	if err := utils.DecodeJson(c, &partialUpdate); err {
		return
	}

	// The code block is checking if the `partialUpdate.Password` is an empty string. If it is empty, it
	// means that the password field was not provided in the request body. In this case, the code block
	// returns a JSON response with a status code of 400 (Bad Request) and an error message indicating that
	// the password is required to update the user.
	if partialUpdate.Password == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "password is required to update the user",
			},
		})
		return
	}

	// The `if !security.ComparePassword(partialUpdate.Password, registeredUser.Password)` condition is
	// checking if the provided password in the request body (`partialUpdate.Password`) matches the stored
	// password for the user (`registeredUser.Password`).
	if !security.ComparePassword(partialUpdate.Password, registeredUser.Password) {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "the password provided is incorrect",
			},
		})
		return
	}

	// The code block is checking if all the fields in the `partialUpdate` struct are empty. If all the
	// fields are empty, it means that no fields were provided in the request body to update the user. In
	// this case, the code block returns a JSON response with a status code of 400 (Bad Request) and an
	// error message indicating that at least one field is required to update the user.
	if partialUpdate.Username == "" && partialUpdate.Email == "" && partialUpdate.FirstName == "" && partialUpdate.LastName == "" && partialUpdate.Age == 0 && partialUpdate.NewPassword == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "at least one field is required to update the user",
			},
		})
		return
	}

	// The code block is checking if the `partialUpdate.Password` is equal to `partialUpdate.NewPassword`.
	// If they are equal, it means that the new password provided in the request body is the same as the
	// old password. In this case, the code block returns a JSON response with a status code of 400 (Bad
	// Request) and an error message indicating that the new password cannot be the same as the old
	// password.
	if partialUpdate.Password == partialUpdate.NewPassword {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "new password cannot be the same as the old password",
			},
		})
		return
	}

	// The code block is checking if the `partialUpdate.NewPassword` field is not empty. If it is not
	// empty, it means that a new password was provided in the request body to update the user's password.
	if partialUpdate.NewPassword != "" {
		hashedPassword, err := security.HashPassword(partialUpdate.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.Response{
				Status: types.Status{
					Code: http.StatusInternalServerError,
					Msg:  "internal server error",
				},
			})
			return
		}
		partialUpdate.NewPassword = hashedPassword
	}

	// The code is creating a new instance of the `User` struct from the `models` package. It assigns the
	// values from the `partialUpdate` struct to the corresponding fields of the `User` struct. This allows
	// for updating the user's information based on the provided values in the request body.
	user := &models.User{
		Username:  partialUpdate.Username,
		Email:     partialUpdate.Email,
		Password:  partialUpdate.NewPassword,
		FirstName: partialUpdate.FirstName,
		LastName:  partialUpdate.LastName,
		Age:       partialUpdate.Age,
	}

	// The code block is calling the `Update` method on the `user` object to update the user's information
	// in the database. The `Update` method takes the ID of the registered user as an argument.
	_, updateErr := user.Update(int(registeredUser.ID))

	if updateErr != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "unable to update the user",
			},
		})
		return
	}

	// The code `c.JSON(http.StatusOK, types.Response{ Status: types.Status{ Code: http.StatusOK, Msg:
	// "user updated successfully", }, })` is returning a JSON response with a status code of 200 (OK) and
	// a message indicating that the user was updated successfully.
	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "user updated successfully",
		},
	})
	return
}

// The `DeleteUser` function is a handler function for deleting a user. It takes a `gin.Context` object
// as a parameter, which provides access to the request and response objects.
func (UserController) DeleteUser(c *gin.Context) {
	// The code block is retrieving the subject (`sub`) and access token (`accessToken`) from the request
	// context using the `getSubFromToken` function. If the `shouldReturn` variable is true, it means that
	// there was an error retrieving the subject and access token, so the function returns and exits the
	// handler function.
	sub, accessToken, shouldReturn := getSubFromToken(c)
	{
		if shouldReturn {
			return
		}
		raw_rt, err := c.Request.Cookie("__rt")
		if err != nil {
			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "something went wrong",
				}})
			return
		}

		cache.RevokeToken(accessToken)
		if raw_rt != nil {
			cache.RevokeToken(raw_rt.Value)
		}

		c.SetCookie("__t", "", -1, "/", "localhost", true, true)
		c.SetCookie("__rt", "", -1, "/", "localhost", true, true)
	}

	// The code block is creating a new instance of the `User` struct from the `models` package.
	var user *models.User

	// The code block is calling the `DeleteByEmail` method on the `user` object to delete the user from
	// the database. The `DeleteByEmail` method takes the subject (`sub`) as an argument.
	if err := user.DeleteByEmail(sub); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "user not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "user deleted successfully",
		},
	},
	)
}

// The function `getAccessToken` retrieves the access token from the Authorization header or a cookie
// in a Gin context.
func getAccessToken(ctx *gin.Context) string {
	var accessToken string
	token := ctx.Request.Header.Get("Authorization")
	if token != "" {
		accessToken = strings.Split(token, " ")[1]
	}
	cookieToken, err := ctx.Request.Cookie("__t")
	if err != nil {
		return ""
	}
	if cookieToken != nil {
		accessToken = cookieToken.Value
	}
	return accessToken
}

// The function `getSub` takes a token as input, verifies it, and returns the subject (sub) from the
// token's claims.
func getSub(token string) (string, error) {
	jwtToken, err := security.VerifyToken(token)
	if err != nil {
		return "", err
	}
	sub, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return sub, nil
}

func getSubFromToken(c *gin.Context) (string, string, bool) {
	// The code block is retrieving the access token from the request context using the `getAccessToken`
	// function. If the access token is empty, it means that it was not provided or could not be
	// retrieved. In that case, a JSON response with a status code of 400 (Bad Request) and an error
	// message is returned.
	accessToken := getAccessToken(c)
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrongl",
			},
		})
		return "", "", true
	}

	// The code block is calling the `getSub` function to retrieve the subject (sub) from the access
	// token. The subject typically represents the user ID or some unique identifier for the user.
	sub, err := getSub(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return "", "", true
	}
	return sub, accessToken, false
}
