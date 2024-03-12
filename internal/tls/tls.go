package tls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/fs"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	tlsFileMode = 0o600
	certPEMFile = "cert.pem"
	keyPEMFile  = "key.pem"
)

var nowFn = time.Now

func LoadClientConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, certPool, err := loadCertificateAndCertPool(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// Deprecated: use LoadClientConfig.
func LoadClientConfigFromDir(dir string) (*tls.Config, error) {
	return LoadClientConfig(filepath.Join(dir, certPEMFile), filepath.Join(dir, keyPEMFile))
}

func LoadServerConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, certPool, err := loadCertificateAndCertPool(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func loadCertificateAndCertPool(certFile, keyFile string) (cert tls.Certificate, _ *x509.CertPool, _ error) {
	certBytes, err := os.ReadFile(certFile)
	if err != nil {
		return cert, nil, errors.WithStack(err)
	}

	privKeyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return cert, nil, errors.WithStack(err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(certBytes) {
		return cert, nil, errors.New("failed to add root certificate to pool")
	}

	cert, err = tls.X509KeyPair(certBytes, privKeyBytes)
	if err != nil {
		return cert, nil, errors.WithStack(err)
	}

	return cert, pool, nil
}

// LoadOrGenerateConfig loads the TLS configuration from the given files,
// or generates a new one if the files do not exist.
func LoadOrGenerateConfig(certFile, keyFile string, logger *zap.Logger) (*tls.Config, error) {
	config, err := LoadServerConfig(certFile, keyFile)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	if config != nil {
		if ttl, err := validateTLSConfig(config); err == nil {
			logger.Info("certificate is valid", zap.Duration("ttl", ttl))
			return config, nil
		}
		logger.Warn("failed to validate TLS config; generating new cartificate", zap.Error(err))
	}

	return generateCertificate(certFile, keyFile)
}

// Deprecated: use LoadOrGenerateConfig.
func LoadOrGenerateConfigFromDir(dir string, logger *zap.Logger) (*tls.Config, error) {
	return LoadOrGenerateConfig(filepath.Join(dir, certPEMFile), filepath.Join(dir, keyPEMFile), logger)
}

func generateCertificate(certFile, keyFile string) (*tls.Config, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "runme",
			Organization: []string{"Stateful, Inc."},
			Country:      []string{"US"},
			Province:     []string{"California"},
			Locality:     []string{"Berkeley"},
		},
		NotBefore:             nowFn(),
		NotAfter:              nowFn().AddDate(0, 0, 30),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IPAddresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
		},
		DNSNames: []string{
			"localhost",
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, err
	}

	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificateBytes,
	}); err != nil {
		return nil, err
	}

	if err := os.WriteFile(certFile, caPEM.Bytes(), tlsFileMode); err != nil {
		return nil, err
	}

	privKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(privKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	}); err != nil {
		return nil, err
	}

	if err := os.WriteFile(keyFile, privKeyPEM.Bytes(), tlsFileMode); err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()

	// TODO: can probably use `AddCert` here
	if !certPool.AppendCertsFromPEM(caPEM.Bytes()) {
		return nil, errors.New("failed to add certificate to certificate pool")
	}

	tlsCA := tls.Certificate{
		Certificate: [][]byte{certificateBytes},
		PrivateKey:  privKey,
		Leaf:        ca,
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCA},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	}

	return tlsConfig, nil
}

func validateTLSConfig(config *tls.Config) (ttl time.Duration, _ error) {
	if len(config.Certificates) < 1 || len(config.Certificates[0].Certificate) < 1 {
		return ttl, errors.New("invalid TLS certificate")
	}

	cert, err := x509.ParseCertificate(config.Certificates[0].Certificate[0])
	if err != nil {
		return ttl, errors.Wrap(err, "failed to parse certificate")
	}

	if nowFn().AddDate(0, 0, 7).After(cert.NotAfter) {
		return ttl, errors.New("certificate will expire soon")
	}

	return nowFn().Sub(cert.NotAfter), nil
}
