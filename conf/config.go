package conf

import (
	"flag"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

// AppConfig contains the application configuration
type AppConfig struct {
	Port        string `ini:"port"`
	Auth        bool   `ini:"auth"`
	TokenMaxAge int    `ini:"token_max_age"`
	CertFile    string `ini:"cert_file"`
	KeyFile     string `ini:"key_file"`
}

// EtcdConfig contains the etcd connection configuration
type EtcdConfig struct {
	RootKey      string `ini:"root_key"`
	DirValue     string `ini:"dir_value"`
	EndPointsRaw string
	EndPoints    []string `ini:"endpoint,omitempty,allowshadow"`
	Auth         bool     `ini:"auth"`
	Username     string   `ini:"username"`
	Password     string   `ini:"password"`
	CertFile     string   `ini:"cert_file"`
	KeyFile      string   `ini:"key_file"`
	CAFile       string   `ini:"ca_file"`
}

// Config contains all configuration options
type Config struct {
	ConfigFile string
	AppConf    AppConfig  `ini:"app"`
	EtcdConf   EtcdConfig `ini:"etcd"`
	PrintVer   bool
}

// Conf is the globaly accessible configuration of the running instance
// TODO: remove global variable
//var Conf Config

// init initializes the command line options for the conf package
func setCMDOptions(config *Config) {
	flag.StringVar(&config.ConfigFile, "configfile", "", "The e3w config file")
	flag.StringVar(&config.AppConf.Port, "port", "8080", "Port bound to web server")
	flag.BoolVar(&config.AppConf.Auth, "auth", false, "Use authentication for web server")
	flag.IntVar(&config.AppConf.TokenMaxAge, "tokenmaxage", 120, "How long is an authentication token valid (in seconds)")
	flag.StringVar(&config.AppConf.CertFile, "certfile", "", "Web server cert file")
	flag.StringVar(&config.AppConf.KeyFile, "keyfile", "", "Web server key file")
	flag.StringVar(&config.EtcdConf.RootKey, "etcdrootkey", "", "Root key (key prefix) used in etcd")
	flag.StringVar(&config.EtcdConf.DirValue, "etcddirvalue", "__etcd_dir_value_fADFbkjqdfs6__", "Value representing directory keys")
	flag.StringVar(&config.EtcdConf.EndPointsRaw, "etcdendpoints", "", "Etcd endpoints (multiple values should be separated by ',')")
	flag.BoolVar(&config.EtcdConf.Auth, "etcdauth", false, "Use authentication for etcd")
	flag.StringVar(&config.EtcdConf.Username, "etcdusername", "", "Username to authenticate against etcd")
	flag.StringVar(&config.EtcdConf.Password, "etcdpassword", "", "Password to authenticate against etcd")
	flag.StringVar(&config.EtcdConf.CertFile, "etcdcertfile", "", "Cert file for authetication against secured etcd endpoint")
	flag.StringVar(&config.EtcdConf.KeyFile, "etcdkeyfile", "", "Key file for authentication against ecured etcd endpoint")
	flag.StringVar(&config.EtcdConf.CAFile, "etcdcafile", "", "CA file (public root cert) for authentication against etcd endpoint")
	flag.BoolVar(&config.PrintVer, "version", false, "Print version")
}

// NewConfig TODO
func NewConfig() (*Config, error) {
	config := &Config{}

	// set command line flags and parse them
	setCMDOptions(config)
	flag.Parse()

	// load config file if provided
	if config.ConfigFile != "" {
		err := InitStructFromINI(&config.AppConf, "app", config.ConfigFile)
		if err != nil {
			return nil, errors.Wrap(err, "InitStructFromINI(app) failed")
		}

		err = InitStructFromINI(&config.EtcdConf, "etcd", config.ConfigFile)
		if err != nil {
			return nil, errors.Wrap(err, "InitStructFromINI(etcd) failed")
		}
	}

	// Reparse flags so that command line options have precedence
	// over the config file
	flag.Parse()

	return config, nil
}

// InitConfig initializes the configuration from a config file
/*func InitConfig() error {
	cfg, err := ini.ShadowLoad(Conf.ConfigFile)
	if err != nil {
		return err
	}

	appSec := cfg.Section("app")
	Conf.Port = appSec.Key("port").Value()
	Conf.Auth = appSec.Key("auth").MustBool()
	Conf.TokenMaxAge = appSec.Key("token_max_age").MustInt()
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
}*/

// InitStructFromINI loads the struct s from the section secName of ini file
// iniURL
func InitStructFromINI(s interface{}, secName, iniURL string) error {
	cfg, err := ini.ShadowLoad(iniURL)
	if err != nil {
		return errors.Wrapf(err, "ini.Load(iniURL) failed: s: %+v; secName: %s; iniURL: %s", s, secName, iniURL)
	}

	sec, err := cfg.GetSection(secName)
	if err != nil {
		return errors.New(fmt.Sprintf("section not found: s: %+v; secName: %s; iniURL: %s", s, secName, iniURL))
	}
	err = sec.MapTo(s)
	if err != nil {
		return errors.New(fmt.Sprintf("MapTo(s) failed: s: %+v; secName: %s; iniURL: %s", s, secName, iniURL))
	}
	log.Printf("DEBUG: InitStructFromINI(): config processed: s: %+v; secName: %s; iniURL: %s", s, secName, iniURL)

	return nil
}
