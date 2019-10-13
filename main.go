package main

import (
	"flag"
	"fmt"
	"github.com/coreos/etcd/version"
	"github.com/gin-gonic/gin"
	"github.com/soyking/e3w/conf"
	"github.com/soyking/e3w/e3ch"
	"github.com/soyking/e3w/routers"
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
	conf.ParseEndPoints()

	fmt.Printf("%+v\n", conf.Conf)

	if conf.Conf.PrintVer {
		fmt.Printf("[%s v%s]\n[etcd %s]\n",
			PROGRAM_NAME, PROGRAM_VERSION,
			version.Version,)
	        os.Exit(0)
	}

	client, err := e3ch.NewE3chClient(&conf.Conf)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.UseRawPath = true
	routers.InitRouters(router, &conf.Conf, client)
	router.Run(":" + conf.Conf.Port)
}
