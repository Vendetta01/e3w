package main

import (
	"fmt"
	"log"
	"os"

	"github.com/VendettA01/e3w/auth"
	"github.com/VendettA01/e3w/conf"
	"github.com/VendettA01/e3w/e3ch"
	"github.com/VendettA01/e3w/routers"
	"github.com/coreos/etcd/version"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
		log.Printf("conf.NetConfig() failed: %+v", err)
		panic(err)
	}

	log.Printf("%+v\n", config)

	if config.PrintVer {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			ProgramName, ProgramVersion,
			version.Version)
		os.Exit(0)
	}

	log.Println("Connecting to etcd...")

	client, err := e3ch.NewE3chClient(config)
	if err != nil {
		log.Printf("ERROR: e3ch.NewE3chClient() failed: %+v", err)
		panic(err)
	}

	userAuths, err := setUpAuth(config)
	if err != nil {
		log.Printf("ERROR: setUpAuths(): error: %+v", err)
		panic(err)
	}

	log.Printf("INFO: main(): userAuths: %#v", userAuths)

	log.Print("INFO: Creating router...")
	router := gin.Default()
	router.UseRawPath = true
	log.Print("INFO: Initializing routers...")
	routers.InitRouters(router, config, client, userAuths)

	if config.AppConf.CertFile != "" && config.AppConf.KeyFile != "" {
		log.Print("INFO: Starting HTTPS server...")
		router.RunTLS(":"+config.AppConf.Port, config.AppConf.CertFile, config.AppConf.KeyFile)
	} else {
		log.Print("INFO: Starting HTTP server...")
		router.Run(":" + config.AppConf.Port)
	}
}
