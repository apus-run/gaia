package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	"github.com/apus-run/gaia/pkg/utils"
)

func makeCA(path string, subject *pkix.Name) (*x509.Certificate, *rsa.PrivateKey, error) {
	// creating a CA which will be used to sign all of our certificates using the x509 package from the Go Standard Library
	caCert := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               *subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10*365, 0, 0),
		IsCA:                  true, // <- indicating this certificate is a CA certificate.
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// generate a private key for the CA
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// create the CA certificate
	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	// Create the CA PEM files
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	if err := os.WriteFile(path+"ca.crt", caPEM.Bytes(), 0644); err != nil {
		return nil, nil, err
	}

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	})
	if err := os.WriteFile(path+"ca.key", caPEM.Bytes(), 0644); err != nil {
		return nil, nil, err
	}
	return caCert, caKey, nil
}

func makeCert(path string, caCert *x509.Certificate, caKey *rsa.PrivateKey, subject *pkix.Name, name string) error {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject:      *subject,
		//IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		return err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(path+name+".crt", certPEM.Bytes(), 0644); err != nil {
		return err
	}

	certKeyPEM := new(bytes.Buffer)
	pem.Encode(certKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	})
	return os.WriteFile(path+name+".key", certKeyPEM.Bytes(), 0644)
}

func TestConfig_Config(t *testing.T) {
	// no TLS
	_tls := TLS{}
	conn, e := _tls.Config()
	assert.Nil(t, conn)
	assert.Nil(t, e)

	// only have insecure option
	_tls = TLS{
		Insecure: true,
	}
	conn, e = _tls.Config()
	assert.NotNil(t, conn)
	assert.Nil(t, e)

	path := utils.GetWorkDir() + "/certs/"
	os.MkdirAll(path, 0755)
	defer os.RemoveAll(path)

	//mTLS
	_tls = TLS{
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
