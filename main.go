package main

import (
	"flag"
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

func setUpAuth() (*auth.UserAuthentications, error) {
	userAuths, err := auth.NewUserAuths()
	if err != nil {
		return nil, errors.Wrap(err, "auth.NewUserAuths() failed")
	}

	if conf.Conf.Auth {
		init := func(userAuth auth.UserAuthentication) error {
			return conf.InitAuthFromINI(userAuth, conf.Conf.ConfigFile)
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
	// Initial parsing of command line options to get ConfigFile
	// if specified
	flag.Parse()
	err := conf.InitConfig()
	if err != nil {
		panic(err)
	}
	// Reparse flags so that command line options have precedence
	// over the config file
	flag.Parse()
	//conf.ParseEndPoints()

	log.Printf("%+v\n", conf.Conf)

	if conf.Conf.PrintVer {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			ProgramName, ProgramVersion,
			version.Version)
		os.Exit(0)
	}

	log.Println("Connecting to etcd...")

	client, err := e3ch.NewE3chClient(&conf.Conf)
	if err != nil {
		log.Printf("ERROR: e3ch.NewE3chClient() failed: %v", err)
		panic(err)
	}

	userAuths, err := setUpAuth()
	if err != nil {
		log.Printf("ERROR: setUpAuths(): error: %+v", err)
		panic(err)
	}

	log.Printf("INFO: main(): userAuths: %#v", userAuths)

	log.Print("INFO: Creating router...")
	router := gin.Default()
	router.UseRawPath = true
	log.Print("INFO: Initializing routers...")
	routers.InitRouters(router, &conf.Conf, client, userAuths)

	if conf.Conf.CertFile != "" && conf.Conf.KeyFile != "" {
		log.Print("INFO: Starting HTTPS server...")
		router.RunTLS(":"+conf.Conf.Port, conf.Conf.CertFile, conf.Conf.KeyFile)
	} else {
		log.Print("INFO: Starting HTTP server...")
		router.Run(":" + conf.Conf.Port)
	}
}
