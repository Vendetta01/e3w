package auth

import (
	"github.com/pkg/errors"
)

// Local defines a struct with the necessary config
type Local struct {
	Username           string `ini:"username"`
	Password           string `ini:"password"`
	AllowEmptyPassword bool   `ini:"allow_empty_password"`
}

// NewLocal returns a new instance of the struct Local
func NewLocal() (*Local, error) {
	return new(Local), nil
}

// login implements the userAuthentication interface method login()
func (l Local) login(userCreds UserCredentials) (bool, error) {
	if userCreds.Username == l.Username &&
		userCreds.Password == l.Password {
		return true, nil
	}
	return false, nil
}

// GetName implements the userAuthentication interface method GetName() (returns "local")
func (l Local) GetName() string {
	return "local"
}

// TestConfig implements the userAuthentication interface method TestConfig()
func (l Local) TestConfig() error {
	if l.Username == "" {
		return errors.New("auth_local: testConfig(): username is empty")
	}
	if l.Password == "" || l.AllowEmptyPassword {
		return errors.New("auth_local: testConfig(): password is empty and not allowed")
	}
	return nil
}
