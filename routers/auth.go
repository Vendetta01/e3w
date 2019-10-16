package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/soyking/e3w/conf"
)

const (
	tokenMaxAge = 120
)

var cache = make(map[string]time.Time)

func getSessionToken() string {
	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	expiresAt := time.Now().Add(tokenMaxAge * time.Second)
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of tokenMaxAge seconds
	cache[sessionToken] = expiresAt

	return sessionToken
}

func authRequired(c *gin.Context) {
	// Check if cookie is present
	fmt.Println("authRequired: checking cookie...")
	userToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("authRequired: no cookie found")
			//c.Redirect(http.StatusSeeOther, "/login")
			/*c.JSON(http.StatusOK, &response{
			Result: nil,
			Err:    errAuthRequired.Error()})*/
			//c.Abort()
			c.AbortWithStatusJSON(http.StatusOK, &response{
				Result: nil,
				Err:    errAuthRequired.Error()})
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Abort()
		return
	}
	fmt.Println("authRequired: cookie successfully read")

	// Check provided session token for validity
	expiresAt, ok := cache[userToken]
	if !ok {
		// Session token is invalid
		fmt.Println("authRequired: session token invalid")
		//c.Redirect(http.StatusSeeOther, "/login")
		/*c.JSON(http.StatusOK, &response{
		Result: nil,
		Err:    errAuthRequired.Error()})*/
		//c.Abort()
		c.AbortWithStatusJSON(http.StatusOK, &response{
			Result: nil,
			Err:    errAuthRequired.Error()})
		return
	}
	if expiresAt.Before(time.Now()) {
		// Session token is expired
		fmt.Println("authRequired: session token expired")
		delete(cache, userToken)
		//c.Redirect(http.StatusSeeOther, "/login")
		/*c.JSON(http.StatusOK, &response{
		Result: nil,
		Err:    errAuthRequired.Error()})*/
		//c.Abort()
		c.AbortWithStatusJSON(http.StatusOK, &response{
			Result: nil,
			Err:    errAuthRequired.Error()})
		return
	}

	// Refresh existing session token
	delete(cache, userToken)
	c.SetCookie("session_token", getSessionToken(), tokenMaxAge, "", "", false, false)
	fmt.Printf("authRequired: Cookie set")

	// Pass on to the next-in-chain
	c.Next()
}

func logIn(c *gin.Context) {
	// First get username and password from POST form
	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Printf("logIn: POST: u: '%v', p: '%v'\nConf: u: '%v', p: '%v'\n", username, password, conf.Conf.Username, conf.Conf.Password)

	if username != conf.Conf.Username || password != conf.Conf.Password {
		// Invalid credentials
		fmt.Println("logIn: username or password missmatch!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte("Unauthorized"))
		return
	}

	c.SetCookie("session_token", getSessionToken(), tokenMaxAge, "", "", false, false)
	c.Redirect(http.StatusSeeOther, "/")
}

func logOut(c *gin.Context) {
	// Check if cookie is present
	userToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.Redirect(http.StatusUnauthorized, "/login")
			return
		}
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, ok := cache[userToken]; ok {
		delete(cache, userToken)
	}

	c.Redirect(http.StatusOK, "/login")
}
