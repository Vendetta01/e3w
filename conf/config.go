package conf

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

type Config struct {
	ConfigFile       string
	Port             string
	Auth             bool
	CertFile         string
	KeyFile          string
	Username         string
	Password         string
	EtcdRootKey      string
	EtcdDirValue     string
	EtcdEndPointsRaw string
	EtcdEndPoints    []string
	EtcdAuth         bool
	EtcdUsername     string
	EtcdPassword     string
	EtcdCertFile     string
	EtcdKeyFile      string
	EtcdCAFile       string
	PrintVer         bool
}

var Conf Config

func init() {
	flag.StringVar(&Conf.ConfigFile, "configfile", "conf/config.default.ini", "The e3w config file")
	flag.StringVar(&Conf.Port, "port", "8080", "Port bound to web server")
	flag.BoolVar(&Conf.Auth, "auth", false, "Use authentication for web server")
	flag.StringVar(&Conf.CertFile, "certfile", "", "Web server cert file")
	flag.StringVar(&Conf.KeyFile, "keyfile", "", "Web server key file")
	flag.StringVar(&Conf.Username, "username", "", "User name for web server auth")
	flag.StringVar(&Conf.Password, "password", "", "Password for web server authentication")
	flag.StringVar(&Conf.EtcdRootKey, "etcdrootkey", "", "Root key (key prefix) used in etcd")
	flag.StringVar(&Conf.EtcdDirValue, "etcddirvalue", "__etcd_dir_value_fADFbkjqdfs6__", "Value representing directory keys")
	flag.StringVar(&Conf.EtcdEndPointsRaw, "etcdendpoints", "", "Etcd endpoints (multiple values should be separated by ',')")
	flag.BoolVar(&Conf.EtcdAuth, "etcdauth", false, "Use authentication for etcd")
	flag.StringVar(&Conf.EtcdUsername, "etcdusername", "", "Username to authenticate against etcd")
	flag.StringVar(&Conf.EtcdPassword, "etcdpassword", "", "Password to authenticate against etcd")
	flag.StringVar(&Conf.EtcdCertFile, "etcdcertfile", "", "Cert file for authetication against secured etcd endpoint")
	flag.StringVar(&Conf.EtcdKeyFile, "etcdkeyfile", "", "Key file for authentication against ecured etcd endpoint")
	flag.StringVar(&Conf.EtcdCAFile, "etcdcafile", "", "CA file (public root cert) for authentication against etcd endpoint")
	flag.BoolVar(&Conf.PrintVer, "version", false, "Print version")
}

func InitConfig() error {
	cfg, err := ini.ShadowLoad(Conf.ConfigFile)
	if err != nil {
		return err
	}

	appSec := cfg.Section("app")
	Conf.Port = appSec.Key("port").Value()
	Conf.Auth = appSec.Key("auth").MustBool()
	Conf.CertFile = appSec.Key("cert_file").Value()
	Conf.KeyFile = appSec.Key("key_file").Value()
	Conf.Username = appSec.Key("username").Value()
	Conf.Password = appSec.Key("password").Value()

	etcdSec := cfg.Section("etcd")
	Conf.EtcdRootKey = etcdSec.Key("root_key").Value()
	Conf.EtcdDirValue = etcdSec.Key("dir_value").Value()
	fmt.Printf("InitConfig: EtcdEndpoints raw: %v\n", etcdSec.Key("addr"))
	fmt.Printf("InitConfig: EtcdEndpoints: %v\n", etcdSec.Key("addr").Value())
	fmt.Printf("InitConfig: EtcdEndpoints with shadows: %v\n", etcdSec.Key("addr").ValueWithShadows())
	Conf.EtcdEndPoints = etcdSec.Key("addr").ValueWithShadows()
	Conf.EtcdAuth = appSec.Key("etcdauth").MustBool()
	Conf.EtcdUsername = etcdSec.Key("username").Value()
	Conf.EtcdPassword = etcdSec.Key("password").Value()
	Conf.EtcdCertFile = etcdSec.Key("cert_file").Value()
	Conf.EtcdKeyFile = etcdSec.Key("key_file").Value()
	Conf.EtcdCAFile = etcdSec.Key("ca_file").Value()

	return nil
}

func ParseEndPoints() {
	Conf.EtcdEndPoints = strings.Split(Conf.EtcdEndPointsRaw, ",")
}
