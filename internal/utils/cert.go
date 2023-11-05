package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/fs"
	"math/big"
	"net"
	"os"
	"time"
)

func CreateCertificate(serverCertPath, privateKeyPath string) error {
	maxInt := 1658
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(int64(maxInt)),
		Subject: pkix.Name{
			Organization: []string{"kripsyInt"},
			Country:      []string{"RU"},
		},
		//nolint:gomnd
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1)},
		NotBefore:   time.Now(),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}
	//nolint:gomnd
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("error generate key %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("error create certificate %w", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return fmt.Errorf("error encode cert %w", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return fmt.Errorf("error encode private key %w", err)
	}

	err = saveCert(serverCertPath, &certPEM)
	if err != nil {
		return err
	}
	err = saveCert(privateKeyPath, &privateKeyPEM)
	if err != nil {
		return err
	}

	return nil
}

func saveCert(path string, payload *bytes.Buffer) error {
	permissionValue := 0755
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(permissionValue))
	if err != nil {
		return fmt.Errorf("error open file %w", err)
	}
	writer := bufio.NewWriter(f)
	_, err = writer.ReadFrom(payload)
	if err != nil {
		return fmt.Errorf("error write to file %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("error close file %w", err)
	}

	return nil
}
