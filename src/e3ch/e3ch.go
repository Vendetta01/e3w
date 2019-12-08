package e3ch

import (
	"crypto/tls"

	"github.com/VendettA01/e3w/src/conf"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	client "github.com/soyking/e3ch"
)

// NewE3chClient returns a new instance of EtcdHRCHYClient
func NewE3chClient(config *conf.Config) (*client.EtcdHRCHYClient, error) {
	var tlsConfig *tls.Config
	var err error
	if config.EtcdConf.CertFile != "" && config.EtcdConf.KeyFile != "" && config.EtcdConf.CAFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      config.EtcdConf.CertFile,
			KeyFile:       config.EtcdConf.KeyFile,
			TrustedCAFile: config.EtcdConf.CAFile,
		}
		tlsConfig, err = tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	clt, err := clientv3.New(clientv3.Config{
		Endpoints: config.EtcdConf.EndPoints,
		Username:  config.EtcdConf.Username,
		Password:  config.EtcdConf.Password,
		TLS:       tlsConfig,
	})
	if err != nil {
		return nil, err
	}

	client, err := client.New(clt, config.EtcdConf.RootKey, config.EtcdConf.DirValue)
	if err != nil {
		return nil, err
	}
	return client, client.FormatRootKey()
}

// CloneE3chClient clones an existing EtcdHRCHYClient (establishes a new connection to
// the etcd server)
func CloneE3chClient(username, password string, client *client.EtcdHRCHYClient) (*client.EtcdHRCHYClient, error) {
	clt, err := clientv3.New(clientv3.Config{
		Endpoints: client.EtcdClient().Endpoints(),
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return client.Clone(clt), nil
}
