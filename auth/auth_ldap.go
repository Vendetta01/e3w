package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/ldap.v3"
)

// Ldap defines a struct with the necessary config
type Ldap struct {
	BindUserDN            string `ini:"bind_user_dn"`
	BindPW                string `ini:"bind_pw"`
	URI                   string `ini:"uri"`
	UseTLS                bool   `ini:"use_tls"`
	TLSInsecureSkipVerify bool   `ini:"tls_insecure_skip_verify"`
	CaCertFile            string `ini:"ca_cert_file"`
	BaseDN                string `ini:"base_dn"`
	UserSearchFilter      string `ini:"user_search_filter"`
}

// NewLdap returns a new instance of the struct Ldap
func NewLdap() (*Ldap, error) {
	return new(Ldap), nil
}

// login implements the userAuthentication interface method login()
func (l Ldap) login(userCreds UserCredentials) (bool, error) {
	con, err := l.connect()
	if err != nil {
		return false, errors.Wrap(err, "auth_ldap: login(): connect() failed")
	}
	defer con.Close()

	// First bind with the specified userdn/pw
	err = con.Bind(l.BindUserDN, l.BindPW)
	if err != nil {
		return false, errors.Wrap(err, "auth_ldap: login(): Bind error")
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		l.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(l.UserSearchFilter, userCreds.Username),
		[]string{"dn"},
		nil,
	)

	searchResult, err := con.Search(searchRequest)
	if err != nil {
		return false, errors.Wrap(err, "auth_ldap: login(): Search() for username failed")
	}

	if len(searchResult.Entries) != 1 {
		return false, errors.New("auth_ldap: login(): user does not exist or search returned multiple results")
	}

	userDN := searchResult.Entries[0].DN

	// Bind as the user to verify their password
	err = con.Bind(userDN, userCreds.Password)
	if err != nil {
		return false, errors.Wrap(err, "auth_ldap: login(): Bind failed")
	}

	// If we get here it means we successfuly bound to the ldap server
	// with the provided user credentials
	return true, nil
}

// GetName implements the userAuthentication interface method GetName() (returns "ldap")
func (l Ldap) GetName() string {
	return "ldap"
}

// TestConfig implements the userAuthentication interface method TestConfig()
func (l Ldap) TestConfig() error {
	con, err := l.connect()
	if err != nil {
		return errors.Wrap(err, "connect() failed")
	}
	defer con.Close()

	// Bind with the specified userdn/pw
	err = con.Bind(l.BindUserDN, l.BindPW)
	if err != nil {
		return errors.Wrap(err, "Bind() failed")
	}

	return nil
}

// connect tries to establish a connection to the ldap server
func (l Ldap) connect() (*ldap.Conn, error) {
	if !l.UseTLS {
		return ldap.Dial("tcp", l.URI)
	}

	rootCA, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "SystemCertPool() failed")
	}
	if rootCA == nil {
		rootCA = x509.NewCertPool()
	}
	ldapCert, err := ioutil.ReadFile(l.CaCertFile)
	if err != nil {
		return nil, errors.Wrap(err, "ReadFile(CaCertFile) failed")
	}
	ok := rootCA.AppendCertsFromPEM(ldapCert)
	if !ok {
		return nil, errors.New("AppendCertsFromPEM() failed")
	}
	tlsConfig := tls.Config{
		InsecureSkipVerify: l.TLSInsecureSkipVerify,
		//ServerName: "",
		RootCAs: rootCA,
	}
	return ldap.DialTLS("tcp", l.URI, &tlsConfig)
}
