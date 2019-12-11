package auth

import (
	"github.com/pkg/errors"
)

// UserCredentials defines a set of credentials (username and password) for a login
// attempt
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//go:generate mockgen -destination mocks/UserAuthentication.go github.com/VendettA01/e3w/src/auth UserAuthentication

// UserAuthentication defines the interface to be implemented by a new authentication method
//
// login(): implements the authentication methods login attempt
//
// GetName(): returns a unique name for the implemented authentication method
//
// TestConfig(): implements a test that returns whether the supplied config can be used
// for login attempts
type UserAuthentication interface {
	login(UserCredentials) (bool, error)
	GetName() string
	TestConfig() error
}

//go:generate mockgen -destination mocks/UserAuthentications.go github.com/VendettA01/e3w/src/auth UserAuthentications

// UserAuthentications defines a collection of authentication methods that are registered
// and provide a valid config
type UserAuthentications struct {
	// Is authentication enabled
	IsEnabled bool
	// AuthMethods contains map of all registered authentication methods
	AuthMethods map[string]UserAuthentication
}

// NewUserAuths returns a new UserAuthentications struct
//
// A number of userAuthentication structs can be provided to initialize
// the returned struct.
//
// TODO: These auth methods are not initialized/have to be initialized before.
// This is contrary to the Register method and should be changed!
func NewUserAuths(userAuths ...UserAuthentication) (*UserAuthentications, error) {
	authMethods := make(map[string]UserAuthentication)
	isEnabled := true

	for _, authMethod := range userAuths {
		_, ok := authMethods[authMethod.GetName()]
		if ok {
			return nil, errAuthNameNotUnique
		}
		authMethods[authMethod.GetName()] = authMethod
	}

	if len(authMethods) < 1 {
		isEnabled = false
	}

	return &UserAuthentications{
		IsEnabled:   isEnabled,
		AuthMethods: authMethods,
	}, nil
}

// RegisterMethod adds an authentication method to the struct UserAuthentications
//
// It expects a user authentication method (UserAuthentication) and an initialization
// function that fills the UserAuthentication struct fields. The function returns
// true or false depending on whether the method was successfuly registered and
// an error code if something went wrong. It is up to the caller to decide how to
// proceed with a faulty registration.
func (userAuths *UserAuthentications) RegisterMethod(userAuth UserAuthentication,
	init func(UserAuthentication) error) (bool, error) {
	_, ok := userAuths.AuthMethods[userAuth.GetName()]
	if ok {
		return false, errAuthNameNotUnique
	}

	// init method
	err := init(userAuth)
	if err != nil {
		return false, errors.Wrap(err, "init() failed")
	}

	userAuths.AuthMethods[userAuth.GetName()] = userAuth
	userAuths.IsEnabled = true

	return true, nil
}

// CanLogIn verifies if the provided credentials are valid
//
// It calls login() on all registered authentication methods and returns true
// as soon as the first one succeeds. If all attempts fail false is returned along
// with the last error code != nil
func (userAuths *UserAuthentications) CanLogIn(userCreds UserCredentials) (bool, error) {
	var lastErr error = nil
	if len(userAuths.AuthMethods) < 1 {
		return false, errNoActiveAuth
	}
	for _, authMethod := range userAuths.AuthMethods {
		authOK, err := authMethod.login(userCreds)
		if err != nil {
			lastErr = err
		}
		if authOK {
			return true, nil
		}
	}
	return false, lastErr
}
