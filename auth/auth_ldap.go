package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/ldap.v3"
)

type AuthLdap struct {
	BindUserDN            string `ini:"binduserdn"`
	BindPW                string `ini:"bindpw"`
	Uri                   string `ini:"uri"`
	UseTLS                bool   `ini:"usetls"`
	TlsInsecureSkipVerify bool   `ini:"tlsinsecureskipverify"`
	CaCertFile            string `ini:"cacertfile"`
	BaseDN                string `ini:"basedn"`
	UserSearchFilter      string `ini:"usersearchfilter"`
}

func (l *AuthLdap) init() bool {
	return initAuthModule(l)
}

func (l AuthLdap) login(userCreds userCredentials) bool {
	con, err := l.connect()
	if err != nil {
		log.Printf("ERROR: auth_ldap: login(): Could not connect: %s", err.Error())
		return false
	}
	defer con.Close()

	// First bind with the specified userdn/pw
	err = con.Bind(l.BindUserDN, l.BindPW)
	if err != nil {
		log.Printf("ERROR: auth_ldap: Could not bind to ldap server: %s", err)
		return false
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
		log.Printf("ERROR: auth_ldap: searchRequest for username failed: %s", err)
		return false
	}

	if len(searchResult.Entries) != 1 {
		log.Printf("ERROR: auth_ldap: login(): User does not exist or too many entries returned")
		return false
	}

	userDN := searchResult.Entries[0].DN

	// Bind as the user to verify their password
	err = con.Bind(userDN, userCreds.Password)
	if err != nil {
		log.Printf("ERROR: auth_ldap: login(): Could not bind with provided user credentials: %s", err)
		return false
	}

	// If we get here it means we successfuly bound to the ldap server
	// with the provided user credentials
	return true
}

func (l AuthLdap) getName() string {
	return "ldap"
}

func (l AuthLdap) testConfig() bool {
	con, err := l.connect()
	if err != nil {
		log.Printf("ERROR: auth_ldap: testConfig(): Could not connect: %s", err.Error())
		return false
	}
	defer con.Close()

	// Bind with the specified userdn/pw
	err = con.Bind(l.BindUserDN, l.BindPW)
	if err != nil {
		log.Printf("ERROR: auth_ldap: testConfig(): Could not bind to ldap server: %s", err)
		return false
	}

	return true
}

func (l AuthLdap) connect() (*ldap.Conn, error) {
	if !l.UseTLS {
		return ldap.Dial("tcp", l.Uri)
	} else {
		rootCA, err := x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		if rootCA == nil {
			log.Printf("INFO: no system root CAs, creating new pool instead")
			rootCA = x509.NewCertPool()
		}
		ldapCert, err := ioutil.ReadFile(l.CaCertFile)
		if err != nil {
			return nil, err
		}
		ok := rootCA.AppendCertsFromPEM(ldapCert)
		if !ok {
			log.Printf("WARN: CaCertFile could not be added: %s", l.CaCertFile)
		}
		tlsConfig := tls.Config{
			InsecureSkipVerify: l.TlsInsecureSkipVerify,
			//ServerName: "",
			RootCAs: rootCA,
		}
		return ldap.DialTLS("tcp", l.Uri, &tlsConfig)
	}
}

func init() {
	authImpl := new(AuthLdap)
	authImpls[authImpl.getName()] = authImpl
}
