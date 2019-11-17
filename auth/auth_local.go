package auth

import (
	"github.com/pkg/errors"
)

// Local TODO
type Local struct {
	Username           string `ini:"username"`
	Password           string `ini:"password"`
	AllowEmptyPassword bool   `ini:"allow_empty_password"`
}

// NewLocal TODO
func NewLocal() (*Local, error) {
	return new(Local), nil
}

func (l Local) login(userCreds UserCredentials) (bool, error) {
	if userCreds.Username == l.Username &&
		userCreds.Password == l.Password {
		return true, nil
	}
	return false, nil
}

// GetName TODO
func (l Local) GetName() string {
	return "local"
}

// TestConfig TODO
func (l Local) TestConfig() error {
	if l.Username == "" {
		return errors.New("auth_local: testConfig(): username is empty")
	}
	if l.Password == "" || l.AllowEmptyPassword {
		return errors.New("auth_local: testConfig(): password is empty and not allowed")
	}
	return nil
}
