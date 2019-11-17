package conf_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/VendettA01/e3w/conf"
)

func TestNewConfigFromINI(t *testing.T) {
	os.Args = []string{"cmd", "-configfile", "test/config.test.ini"}
	config, err := conf.NewConfig()

	if err != nil {
		t.Fatal(err)
	}
	// test if auth was read correctly (default is false)
	if config.AppConf.Auth != true {
		t.Error("[app].auth != true")
	}
	// test if multiple endpoints are read
	exp := []string{"etcd:2379", "etcd:22379", "etcd:32379"}
	act := config.EtcdConf.EndPoints
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("expected: %v, actual: %v", exp, act)
	}
}
