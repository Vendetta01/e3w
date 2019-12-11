package routers

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

// constants TODO
const (
	ETCDClientTimeout = 3 * time.Second
)

func parseBody(c *gin.Context, t interface{}) error {
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, t)
}

func newEtcdCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), ETCDClientTimeout)
	return ctx
}
