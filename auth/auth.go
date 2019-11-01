package auth

import (
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	"time"

	"github.com/VendettA01/e3w/conf"
	"github.com/VendettA01/e3w/resp"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

type userCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var cache = make(map[string]time.Time)

func getSessionToken() string {
	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	expiresAt := time.Now().Add(time.Duration(conf.Conf.TokenMaxAge) * time.Second)
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of tokenMaxAge seconds
	cache[sessionToken] = expiresAt

	return sessionToken
}

func AuthRequired(c *gin.Context) {
	// if authentication is disabled continue with next handler
	if !conf.Conf.Auth {
		c.Next()
		return
	}

	// Check if cookie is present
	fmt.Println("authRequired: checking cookie...")
	userToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("authRequired: no cookie found")
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
	fmt.Println("authRequired: cookie successfully read")

	// Check provided session token for validity
	expiresAt, ok := cache[userToken]
	if !ok {
		// Session token is invalid
		fmt.Println("authRequired: session token invalid")
		c.AbortWithStatusJSON(http.StatusUnauthorized, &resp.Response{
			Result: nil,
			Err:    errAuthRequired.Error()})
		return
	}
	if expiresAt.Before(time.Now()) {
		// Session token is expired
		fmt.Println("authRequired: session token expired")
		delete(cache, userToken)
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			&resp.Response{
				Result: nil,
				Err:    errAuthRequired.Error()})
		return
	}

	// Refresh existing session token
	delete(cache, userToken)
	c.SetCookie("session_token", getSessionToken(), conf.Conf.TokenMaxAge, "", "", false, false)
	fmt.Println("authRequired: Cookie set")

	// Pass on to the next-in-chain
	c.Next()
}

func LogIn(c *gin.Context) {
	// First get username and password from POST form
	var userCreds userCredentials
	err := json.NewDecoder(c.Request.Body).Decode(&userCreds)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &resp.Response{
			Result: nil,
			Err:    errInvJSONOnRequest.Error()})
		return
	}

	fmt.Printf("logIn: POST: u: '%v', p: '%v'\n", userCreds.Username, userCreds.Password)

	loginSucessful, err := canLogIn(userCreds)
	if err != nil {
		// Some internal error occured, pass it on
		c.AbortWithStatusJSON(http.StatusUnauthorized, &resp.Response{
			Result: nil,
			Err:    err.Error()})
		return
	}
	if !loginSucessful {
		// Invalid credentials
		fmt.Println("logIn: username or password missmatch!")
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			&resp.Response{
				Result: nil,
				Err:    errInvCredentials.Error()})
		return
	}

	c.SetCookie("session_token", getSessionToken(), conf.Conf.TokenMaxAge, "", "", false, false)
}

func CheckToken(c *gin.Context) {
	// The way the route "/checkToken" is designed, the validity will
	// be checked before this handler is called. If we arrive here
	// it means that authentication was successful
	c.JSON(http.StatusOK, nil)
}

func LogOut(c *gin.Context) {
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
