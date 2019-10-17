package main

import (
	"flag"
	"fmt"
	"github.com/VendettA01/e3w/conf"
	"github.com/VendettA01/e3w/e3ch"
	"github.com/VendettA01/e3w/routers"
	"github.com/coreos/etcd/version"
	"github.com/gin-gonic/gin"
	"os"
)

const (
	PROGRAM_NAME    = "e3w"
	PROGRAM_VERSION = "0.0.3"
)

/*func init() {
	flag.StringVar(&configFilepath, "conf", "conf/config.default.ini", "config file path")
	rev := flag.Bool("rev", false, "print rev")
	flag.Parse()

	if *rev {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version,
		)
		os.Exit(0)
	}
}*/

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

	fmt.Printf("%+v\n", conf.Conf)

	if conf.Conf.PrintVer {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version)
		os.Exit(0)
	}

	fmt.Println("Connecting to etcd...")

	client, err := e3ch.NewE3chClient(&conf.Conf)
	if err != nil {
		fmt.Println("ERROR")
		panic(err)
	}

	fmt.Println("Creating router...")
	router := gin.Default()
	router.UseRawPath = true
	fmt.Println("Initializing routers...")
	routers.InitRouters(router, &conf.Conf, client)

	if conf.Conf.CertFile != "" && conf.Conf.KeyFile != "" {
		fmt.Println("Starting HTTPS server...")
		router.RunTLS(":"+conf.Conf.Port, conf.Conf.CertFile, conf.Conf.KeyFile)
	} else {
		fmt.Println("Starting HTTP server...")
		router.Run(":" + conf.Conf.Port)
	}
}
