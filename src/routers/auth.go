package routers

import (
	"encoding/json"
	"fmt"

	"net/http"
	"time"

	"github.com/VendettA01/e3w/src/auth"
	"github.com/VendettA01/e3w/src/conf"
	"github.com/VendettA01/e3w/src/resp"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"

	log "github.com/sirupsen/logrus"
)

// this is the global token cache
// TODO: should be abstracted so that you can use anything as a backend
var cache = make(map[string]time.Time)

func getSessionToken(tokenMaxAge int) string {
	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	expiresAt := time.Now().Add(time.Duration(tokenMaxAge) * time.Second)
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of tokenMaxAge seconds
	cache[sessionToken] = expiresAt

	return sessionToken
}

// authRequired TODO
func authRequired(userAuths *auth.UserAuthentications, config *conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if cookie is present
		log.Debug("authRequired: checking cookie...")
		userToken, err := c.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Debug("authRequired: no cookie found")
				c.AbortWithStatusJSON(http.StatusUnauthorized,
					&resp.Response{
						Result: nil,
						Err:    errAuthRequired.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest,
				&resp.Response{
					Result: nil,
					Err:    errAuthRequired.Error()})
			return
		}
		log.Debug("authRequired: cookie successfully read")

		// Check provided session token for validity
		expiresAt, ok := cache[userToken]
		if !ok {
			// Session token is invalid
			log.Debug("authRequired: session token invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, &resp.Response{
				Result: nil,
				Err:    errAuthRequired.Error()})
			return
		}
		if expiresAt.Before(time.Now()) {
			// Session token is expired
			log.Debug("authRequired: session token expired")
			delete(cache, userToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				&resp.Response{
					Result: nil,
					Err:    errAuthRequired.Error()})
			return
		}

		// Refresh existing session token
		delete(cache, userToken)
		tokenMaxAge := config.AppConf.TokenMaxAge
		c.SetCookie("session_token", getSessionToken(tokenMaxAge), tokenMaxAge, "", "", false, false)
		log.Debug("authRequired: Cookie set")
	}
}

// logIn TODO
func logIn(userAuths *auth.UserAuthentications, config *conf.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First get username and password from POST form
		var userCreds auth.UserCredentials
		err := json.NewDecoder(c.Request.Body).Decode(&userCreds)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &resp.Response{
				Result: nil,
				Err:    errInvJSONOnRequest.Error()})
			return
		}

		log.WithField("userCreds", fmt.Sprintf("%+v", userCreds)).Debug("logIn: POST")

		loginSucessful, err := userAuths.CanLogIn(userCreds)
		if err != nil {
			// Some internal error occured, pass it on
			c.AbortWithStatusJSON(http.StatusUnauthorized, &resp.Response{
				Result: nil,
				Err:    err.Error()})
			return
		}
		if !loginSucessful {
			// Invalid credentials
			log.Debug("logIn: username or password missmatch!")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				&resp.Response{
					Result: nil,
					Err:    errInvCredentials.Error()})
			return
		}

		tokenMaxAge := config.AppConf.TokenMaxAge
		c.SetCookie("session_token", getSessionToken(tokenMaxAge), tokenMaxAge, "", "", false, false)
	}
}

// checkToken TODO
func checkToken(c *gin.Context) {
	// The way the route "/checkToken" is designed, the validity will
	// be checked before this handler is called. If we arrive here
	// it means that authentication was successful
	c.JSON(http.StatusOK, nil)
}

// logOut TODO
func logOut(c *gin.Context) {
	// Check if cookie is present
	userToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				&resp.Response{
					Result: nil,
					Err:    errAuthRequired.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest,
			&resp.Response{
				Result: nil,
				Err:    errAuthRequired.Error()})
		return
	}

	if _, ok := cache[userToken]; ok {
		delete(cache, userToken)
	}

	c.JSON(http.StatusOK, nil)
}
