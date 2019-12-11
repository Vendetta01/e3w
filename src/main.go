package main

import (
	"fmt"
	"os"

	"github.com/VendettA01/e3w/src/auth"
	"github.com/VendettA01/e3w/src/conf"
	"github.com/VendettA01/e3w/src/e3ch"
	"github.com/VendettA01/e3w/src/routers"
	"github.com/coreos/etcd/version"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

// program constants TODO
const (
	ProgramName    = "e3w"
	ProgramVersion = "0.1.0"
)

func setUpAuth(config *conf.Config) (*auth.UserAuthentications, error) {
	userAuths, err := auth.NewUserAuths()
	if err != nil {
		return nil, errors.Wrap(err, "auth.NewUserAuths() failed")
	}

	if config.AppConf.Auth {
		init := func(userAuth auth.UserAuthentication) error {
			err := conf.InitStructFromINI(userAuth, "auth:"+userAuth.GetName(), config.ConfigFile)
			if err != nil {
				return err
			}
			return userAuth.TestConfig()
		}
		authLocal, err := auth.NewLocal()
		if err != nil {
			return nil, errors.Wrap(err, "auth.NewLocal() failed")
		}
		ok, err := userAuths.RegisterMethod(authLocal, init)
		if !ok {
			log.Printf("WARN: setUpAuth(): auth_local: not registered: %s", err)
		}

		authLdap, err := auth.NewLdap()
		if err != nil {
			return nil, errors.Wrap(err, "auth.NewLdap() failed")
		}
		ok, err = userAuths.RegisterMethod(authLdap, init)
		if !ok {
			log.Printf("WARN: setUpAuth(): auth_ldap: not registered: %s", err)
		}
	}

	return userAuths, nil
}

func main() {
	config, err := conf.NewConfig()
	if err != nil {
		log.WithField("err", fmt.Sprintf("%+v", err)).Error("conf.NewConfig() failed")
	}

	// enable debugging if set
	if config.AppConf.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.WithField(
		"config", fmt.Sprintf("%+v", config)).Debug("Configuration successfuly parsed")

	if config.PrintVer {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			ProgramName, ProgramVersion,
			version.Version)
		os.Exit(0)
	}

	log.Info("Connecting to etcd...")

	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		log.WithField("err", fmt.Sprintf("%+v", err)).Error("NewE3chClient() failed")
	}

	userAuths, err := setUpAuth(config)
	if err != nil {
		log.WithField("err", fmt.Sprintf("%+v", err)).Error("setUpAuth() failed")
	}

	if config.AppConf.Auth && !userAuths.IsEnabled {
		log.Error("No user auth method successfuly registered but auth is enabled in config")
	}

	log.WithField("userAuths",
		fmt.Sprintf("%+v", userAuths)).Debug("Authentication methods registered")

	log.Debug("Creating router...")
	router := gin.Default()
	router.UseRawPath = true
	log.Info("Initializing routers...")
	routers.InitRouters(router, config, client, userAuths)

	if config.AppConf.CertFile != "" && config.AppConf.KeyFile != "" {
		log.Info("Starting HTTPS server...")
		router.RunTLS(":"+config.AppConf.Port, config.AppConf.CertFile, config.AppConf.KeyFile)
	} else {
		log.Info("Starting HTTP server...")
		router.Run(":" + config.AppConf.Port)
	}
}
