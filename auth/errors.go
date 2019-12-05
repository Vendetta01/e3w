package auth

import "github.com/pkg/errors"

// Errors
var (
	errNoActiveAuth      = errors.New("no active auth")
	errAuthNameNotUnique = errors.New("authentication method names are not unique")
)
