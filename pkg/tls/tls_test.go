package tls

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Config(t *testing.T) {
	// no TLS
	_tls := Config{}
	conn, e := _tls.Config()
	assert.Nil(t, conn)
	assert.Nil(t, e)

	// only have insecure option
	_tls = Config{
		Insecure: true,
	}
	conn, e = _tls.Config()
	assert.NotNil(t, conn)
	assert.Nil(t, e)

	path := GetWorkDir() + "/certs/"
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	//mTLS
	_tls = Config{
		CA:   filepath.Join(path, "./ca.crt"),
		Cert: filepath.Join(path, "./test.crt"),
		Key:  filepath.Join(path, "./test.key"),
	}
	conn, e = _tls.Config()
	assert.Nil(t, conn)
	assert.NotNil(t, e)

	subject := pkix.Name{
		Country:            []string{"Earth"},
		Organization:       []string{"MegaEase"},
		OrganizationalUnit: []string{"Engineering"},
		Locality:           []string{"Mountain"},
		Province:           []string{"Asia"},
		StreetAddress:      []string{"Bridge"},
		PostalCode:         []string{"123456"},
		SerialNumber:       "",
		CommonName:         "CA",
		Names:              []pkix.AttributeTypeAndValue{},
		ExtraNames:         []pkix.AttributeTypeAndValue{},
	}
	caCert, caKey, err := makeCA(path, &subject)
	if err != nil {
		t.Fatalf("make CA Certificate error! - %v", err)
	}
	t.Log("Create the CA certificate successfully.")

	subject.CommonName = "Server"
	subject.Organization = []string{"Server Company"}
	if err := makeCert(path, caCert, caKey, &subject, "test"); err != nil {
		t.Fatal("make Server Certificate error!")
	}
	t.Log("Create and Sign the Server certificate successfully.")

	conn, e = _tls.Config()
	assert.Nil(t, e)
	assert.NotNil(t, conn)

	monkey.Patch(tls.LoadX509KeyPair, func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, fmt.Errorf("load x509 key pair error")
	})

	conn, e = _tls.Config()
	assert.NotNil(t, e)
	assert.Nil(t, conn)
	monkey.UnpatchAll()

	//TLS
	_tls = TLS{
		CA:       filepath.Join(path, "./ca.crt"),
		Insecure: false,
	}
	conn, e = _tls.Config()
	assert.Nil(t, e)
	assert.NotNil(t, conn)
	assert.Nil(t, conn.Certificates)
}
