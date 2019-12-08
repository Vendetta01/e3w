package conf

import (
	"flag"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

// stringSliceFlag defines a new type (string slice) for processing the same flag
// multiple times
type stringSliceFlag []string

// String implements the String() method of the Value interface of package flag
func (ssf *stringSliceFlag) String() string {
	ret := "["
	for i, s := range *ssf {
		if i != 0 {
			ret += ", "
		}
		ret += s
	}
	ret += "]"
	return ret
}

// Set implements the Set() method of the Value interface of package flag
func (ssf *stringSliceFlag) Set(value string) error {
	*ssf = append(*ssf, value)
	return nil
}

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
	RootKey   string          `ini:"root_key"`
	DirValue  string          `ini:"dir_value"`
	EndPoints stringSliceFlag `ini:"endpoint,omitempty,allowshadow"`
	Auth      bool            `ini:"auth"`
	Username  string          `ini:"username"`
	Password  string          `ini:"password"`
	CertFile  string          `ini:"cert_file"`
	KeyFile   string          `ini:"key_file"`
	CAFile    string          `ini:"ca_file"`
}

// Config contains all configuration options
type Config struct {
	ConfigFile string
	AppConf    AppConfig  `ini:"app"`
	EtcdConf   EtcdConfig `ini:"etcd"`
	PrintVer   bool
}

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
	flag.Var(&config.EtcdConf.EndPoints, "etcdendpoint", "Etcd endpoint (parameter can be set multiple times)")
	flag.BoolVar(&config.EtcdConf.Auth, "etcdauth", false, "Use authentication for etcd")
	flag.StringVar(&config.EtcdConf.Username, "etcdusername", "", "Username to authenticate against etcd")
	flag.StringVar(&config.EtcdConf.Password, "etcdpassword", "", "Password to authenticate against etcd")
	flag.StringVar(&config.EtcdConf.CertFile, "etcdcertfile", "", "Cert file for authetication against secured etcd endpoint")
	flag.StringVar(&config.EtcdConf.KeyFile, "etcdkeyfile", "", "Key file for authentication against ecured etcd endpoint")
	flag.StringVar(&config.EtcdConf.CAFile, "etcdcafile", "", "CA file (public root cert) for authentication against etcd endpoint")
	flag.BoolVar(&config.PrintVer, "version", false, "Print version")
}

// NewConfig returns a new instance of the struct Config
func NewConfig() (*Config, error) {
	config := &Config{}

	// set command line flags and parse them
	setCMDOptions(config)
	flag.Parse()

	clearEndpoints := false
	if len(config.EtcdConf.EndPoints) > 0 {
		clearEndpoints = true
	}

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
	// clear EndPoints if they exist in flags since otherwise we would have
	// the union of config file and flags
	if clearEndpoints {
		config.EtcdConf.EndPoints = []string{}
	}
	flag.Parse()

	return config, nil
}

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
