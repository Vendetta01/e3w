package auth

import (
	"log"
	//"github.com/VendettA01/e3w/conf"
	//"gopkg.in/ini.v1"
)

type AuthLocal struct {
	Username string `ini:"username"`
	Password string `ini:"password"`
}

func (l *AuthLocal) init() bool {
	return initAuthModule(l)
}

func (l AuthLocal) login(userCreds userCredentials) bool {
	if userCreds.Username == l.Username &&
		userCreds.Password == l.Password {
		return true
	}
	return false
}

func (l AuthLocal) getName() string {
	return "local"
}

func (l AuthLocal) testConfig() bool {
	if l.Username == "" {
		log.Printf("ERROR: auth_local: testConfig(): username is empty")
		return false
	}
	if l.Password == "" {
		log.Printf("WARN: auth_local: password is empty, this is not recommended")
	}
	return true
}

func init() {
	authImpl := new(AuthLocal)
	authImpls[authImpl.getName()] = authImpl
}
