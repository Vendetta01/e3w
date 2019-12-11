package conf

import (
	"os"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewConfigFromINI(t *testing.T) {
	os.Args = []string{"cmd", "-configfile", "test/config.test.ini"}
	config, err := NewConfig()

	if err != nil {
		t.Fatal(err)
	}
	// test if auth was read correctly (default is false)
	if config.AppConf.Auth != true {
		t.Error("[app].auth != true")
	}
	// test if multiple endpoints are read
	exp := stringSliceFlag{"etcd:2379", "etcd:22379", "etcd:32379"}
	act := config.EtcdConf.EndPoints
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("expected: %v, actual: %v", exp, act)
	}
}

func TestNewConfigFromCMD(t *testing.T) {
	os.Args = []string{"cmd", "-version",
		"-configfile=test/config.test.ini",
		"-etcdendpoint=testendpoint1",
		"-etcdendpoint=testendpoint2",
		"-auth=false"}
	config, err := NewConfig()

	log.WithField("config", config).Info("DEBUG: config:")

	if err != nil {
		t.Fatal(err)
	}

	if !config.PrintVer {
		t.Error("config.PrintVer: exp: true; act: false")
	}

	exp := stringSliceFlag{"testendpoint1", "testendpoint2"}
	act := config.EtcdConf.EndPoints
	if !reflect.DeepEqual(act, exp) {
		t.Errorf("config.EtcdConf.Endpoints: act: %+v; exp: %+v", act, exp)
	}
}
