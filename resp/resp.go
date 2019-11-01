package resp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Result interface{} `json:"result"`
	Err    string      `json:"err"`
}

type RespHandler func(c *gin.Context) (interface{}, error)

func Resp(handler RespHandler) gin.HandlerFunc {
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
