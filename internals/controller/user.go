package controller

import (
	"errors"
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

type UpdateUser struct {
	Username    string `json:"username,omitempty" validate:"omitempty,min=3,max=32,alphanum"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	Password    string `json:"password,omitempty" validate:"required,min=8"`
	NewPassword string `json:"new_password,omitempty" validate:"omitempty,min=8"`
	FirstName   string `json:"firstname,omitempty" validate:"omitempty,alpha"`
	LastName    string `json:"lastname,omitempty" validate:"omitempty,alpha"`
	Age         int    `json:"age,omitempty" validate:"omitempty,gt=0,lt=100"`
}

func (u *UserController) Get(c *gin.Context) {
	var accessToken string

	haveErr := extractAndValidateToken(c, &accessToken)
	if haveErr {
		return
	}

	email, err := extractEmailFromToken(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return
	}

	var user models.User
	if err := user.GetUserByEmail(email); err != nil {
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
		Data: user,
	})

}

func (u *UserController) Update(c *gin.Context) {
	var update UpdateUser
	var accessToken string

	haveErr := extractAndValidateToken(c, &accessToken)
	if haveErr {
		return
	}

	email, err := extractEmailFromToken(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return
	}

	var user models.User
	if err := user.GetUserByEmail(email); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "user not found",
			},
		})
		return
	}

	isJsonDecoded := utils.DecodeJson(c, &update)
	if isJsonDecoded {
		return
	}

	if err := validate.Struct(&update); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "validation error",
			},
			Errors: utils.ConvertValidationErrors(err),
		})
		return
	}

	if update.Username != "" && update.Email != "" && update.NewPassword != "" && update.FirstName != "" && update.LastName != "" && update.Age != 0 {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "at least one field is required",
			},
		})
		return
	}

	if !security.ComparePassword(update.Password, user.Password) {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "password is incorrect",
			},
		})
		return
	}

	if update.Username != "" {
		var check models.User
		if err := check.GetUserByUsername(update.Username); err == nil {
			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "username already exists",
				},
			})
			return
		}
	}

	if update.Email != "" {
		var check models.User
		if err := check.GetUserByEmail(update.Email); err == nil {
			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "email already exists",
				},
			})
			return
		}
	}

	if update.NewPassword != "" {
		update.NewPassword, err = security.HashPassword(update.NewPassword)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  "something went wrong",
				},
			})
			return
		}
	}

	updateUser := models.User{
		Username:  update.Username,
		Email:     update.Email,
		Password:  update.NewPassword,
		FirstName: update.FirstName,
		LastName:  update.LastName,
		Age:       update.Age,
	}

	if _, err := updateUser.Update(int(user.ID)); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "ok",
		},
		Data: updateUser,
	})
}

func (u *UserController) Delete(c *gin.Context) {
	var accessToken string

	haveErr := extractAndValidateToken(c, &accessToken)
	if haveErr {
		return
	}

	email, err := extractEmailFromToken(accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return
	}

	var user models.User
	if err := user.GetUserByEmail(email); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "user not found",
			},
		})
		return
	}
	var deleteUser models.User
	if err := deleteUser.Delete(int(user.ID)); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return
	}

	token := c.Request.Header.Get("Authorization")
	if token != "" {
		accessToken := strings.Split(token, " ")[1]
		cache.RevokeToken(accessToken)
	}
	c.SetCookie("__t", "", -1, "/", "localhost", true, true)
	c.SetCookie("__rt", "", -1, "/", "localhost", true, true)

	c.JSON(http.StatusOK, types.Response{
		Status: types.Status{
			Code: http.StatusOK,
			Msg:  "ok",
		},
	})
}

func extractEmailFromToken(token string) (string, error) {
	jwtToken, err := security.VerifyToken(token)

	if err != nil {
		return "", err
	}

	if jwtToken.Valid {
		claims, err := jwtToken.Claims.GetSubject()
		if err != nil {
			return "", err
		}
		return claims, nil
	}

	return "", errors.New("something went wrong")
}

func extractAndValidateToken(c *gin.Context, accessToken *string) bool {
	headerToken := c.Request.Header.Get("Authorization")
	cookieToken, err := c.Request.Cookie("__t")

	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  "something went wrong",
			},
		})
		return true
	}

	if headerToken == "" && cookieToken == nil {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status: types.Status{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
			},
		})
		return true
	}

	if headerToken != "" {
		*accessToken = headerToken
	}
	if cookieToken != nil {
		*accessToken = cookieToken.Value
	}

	if security.IsTokenExpired(*accessToken) && cache.IsTokenRevoked(*accessToken) {
		c.JSON(http.StatusUnauthorized, types.Response{
			Status: types.Status{
				Code: http.StatusUnauthorized,
				Msg:  "Unauthorized",
			},
		})
		return true
	}
	return false
}
