package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response defines the structure of a response of the web server
type Response struct {
	Result interface{} `json:"result"`
	Err    string      `json:"err"`
}

// HandlerType defines a handler wrapper
type HandlerType func(c *gin.Context) (interface{}, error)

// Resp implements the response handler
func Resp(handler HandlerType) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		r := &Response{}
		if err != nil {
			r.Err = err.Error()
		} else {
			r.Result = result
		}
		c.JSON(http.StatusOK, r)
	}
}
