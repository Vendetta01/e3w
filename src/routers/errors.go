package routers

import "errors"

var (
	errRoleName        = errors.New("role's name should not be empty")
	errInvalidPermType = errors.New("perm type should be READ | WRITE | READWRITE")

	errUserName         = errors.New("user's name should not be empty")
	errAuthRequired     = errors.New("authentication required")
	errInvJSONOnRequest = errors.New("invalid JSON on request")
	errInvCredentials   = errors.New("invalid credentials")
)
