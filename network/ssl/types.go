package ssl

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/fengzhongzhu1621/xgo/config/viper/viper_parser"
)

type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	InsecureSkipVerify bool
	// Server requires TLS client certificate authentication
	CertFile string
	// Server requires TLS client certificate authentication
	KeyFile string
	// Trusted root certificates for server
	CAFile string
	// the password to decrypt the certificate
	Password string
}

// NewTLSClientConfigFromConfig new config about tls client config
func NewTLSClientConfigFromConfig(prefix string) (TLSClientConfig, error) {
	tlsConfig := TLSClientConfig{}

	skipVerifyKey := fmt.Sprintf("%s.insecureSkipVerify", prefix)
	if val, err := viper_parser.String(skipVerifyKey); err == nil {
		skipVerifyVal := val
		if skipVerifyVal == "true" {
			tlsConfig.InsecureSkipVerify = true
		}
	}

	certFileKey := fmt.Sprintf("%s.certFile", prefix)
	if val, err := viper_parser.String(certFileKey); err == nil {
		tlsConfig.CertFile = val
	}

	keyFileKey := fmt.Sprintf("%s.keyFile", prefix)
	if val, err := viper_parser.String(keyFileKey); err == nil {
		tlsConfig.KeyFile = val
	}

	caFileKey := fmt.Sprintf("%s.caFile", prefix)
	if val, err := viper_parser.String(caFileKey); err == nil {
		tlsConfig.CAFile = val
	}

	passwordKey := fmt.Sprintf("%s.password", prefix)
	if val, err := viper_parser.String(passwordKey); err == nil {
		tlsConfig.Password = val
	}

	return tlsConfig, nil
}

// ExtraClientConfig extra http client configuration
type ExtraClientConfig struct {
	// ResponseHeaderTimeout the amount of time to wait for a server's response headers
	ResponseHeaderTimeout time.Duration
}

// GetClientTLSConfig get client tls config
func GetClientTLSConfig(prefix string) (*tls.Config, error) {
	config, err := NewTLSClientConfigFromConfig(prefix)
	if err != nil {
		return nil, err
	}
	tlsConf := &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}

	if len(config.CAFile) != 0 && len(config.CertFile) != 0 && len(config.KeyFile) != 0 {
		tlsConf, err = ClientTslConfVerity(config.CAFile, config.CertFile, config.KeyFile, config.Password)
		if err != nil {
			return nil, err
		}
	}

	return tlsConf, nil
}
