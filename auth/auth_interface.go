package auth

import (
	"fmt"
	"log"

	"github.com/VendettA01/e3w/conf"
	"gopkg.in/ini.v1"
)

type userAuthentication interface {
	init() bool
	login(userCredentials) bool
	getName() string
	testConfig() bool
}

var authImpls = make(map[string]userAuthentication)

var activeAuths = make([]string, 0, 2)

func InitAuthFromConf() error {
	// if authentication is disabled there is nothing to set up
	if !conf.Conf.Auth {
		return nil
	}
	for authName, authImpl := range authImpls {
		if initSuccess := authImpl.init(); initSuccess {
			activeAuths = append(activeAuths, authName)
			log.Printf("INFO: auth_interface: InitAuthFromConf(): '%s' successfuly registered", authName)
		}
	}

	if len(activeAuths) < 1 {
		// no active authentication method found
		return errNoActiveAuth
	}

	/*log.Printf("DEBUG: auth_interface: InitAuthFromConf: authImpls: %+v", authImpls)
	log.Printf("DEBUG: auth_interface: InitAuthFromConf: activeAuths: %+v", activeAuths)
	for i, authName := range activeAuths {
		log.Printf("DEBUG: auth_interface:InitAuthFromConf: activeAuths: %+v: authImpls[]: %+v: i: %+v", authName, authImpls[authName], i)
	}*/

	return nil
}

func initAuthModule(userAuth userAuthentication) bool {
	cfg, err := ini.Load(conf.Conf.ConfigFile)
	if err != nil {
		log.Printf("ERROR: Cannot open config file for auth module: %s", userAuth.getName())
		return false
	}

	cfg.Section(fmt.Sprintf("auth:%s", userAuth.getName())).MapTo(userAuth)
	log.Printf("DEBUG: initAuthModule(): config processed for: %s: %+v", userAuth.getName(), userAuth)

	return userAuth.testConfig()
}

func canLogIn(userCreds userCredentials) (bool, error) {
	log.Printf("DEBUG: canLogIn(): userCreds: %+v\n", userCreds)
	for _, authName := range activeAuths {
		if authImpl, ok := authImpls[authName]; ok {
			return authImpl.login(userCreds), nil
		} else {
			log.Panic("auth: activeAuth not found in authImpls")
		}
	}
	return false, errNoActiveAuth
}
