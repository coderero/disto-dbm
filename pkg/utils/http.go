package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"coderero.dev/projects/go/gin/hello/types"
	"github.com/gin-gonic/gin"
)

// The function checks if the content type of the request is the expected type and returns true if it
// is not.
func CheckContentType(ctx *gin.Context, t string) bool {
	if ctx.Request.Header.Get("Content-Type") != t {
		ctx.JSON(422, types.Response{
			Status: types.Status{
				Code: http.StatusUnprocessableEntity,
				Msg:  fmt.Sprintf("Content-Type must be %s", t),
			},
		})
		return true
	}
	return false
}

// The function `decodeJson` decodes a JSON request body into a given model, handles any decoding or
// validation errors, and returns a boolean indicating whether an error occurred.
// NOTE: Although Gin has a built-in JSON decoder, but it is not as flexible to decode JSON request
// body into a given model and handle any decoding or validation errors.
func DecodeJson(c *gin.Context, model interface{}) bool {
	// The code below is decoding the request body into the model provided and checking for any errors in
	// the process.
	if err := json.NewDecoder(c.Request.Body).Decode(model); err != nil {
		if strings.Contains(err.Error(), "json:") {
			c.JSON(http.StatusBadRequest, types.Response{
				Status: types.Status{
					Code: http.StatusBadRequest,
					Msg:  ExtractInformation(err),
				},
			})
			return true
		}
		c.JSON(http.StatusBadRequest, types.Response{
			Status: types.Status{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			},
		})
		return true
	}
	return false
}
