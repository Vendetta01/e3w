package auth

import "errors"

var (
	errUserName               = errors.New("auth: user's name should not be empty")
	errAuthRequired           = errors.New("auth: authentication required")
	errInvJSONOnRequest       = errors.New("auth: invalid JSON on request")
	errInvCredentials         = errors.New("auth: invalid credentials")
	errAuthLocalUsernameEmpty = errors.New("auth: local: username is empty")
	errNoActiveAuth           = errors.New("auth: no active auth")
)
