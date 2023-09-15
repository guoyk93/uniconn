package uniconn

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"
)

// DialConfig is the configuration for dialing a connection.
type DialConfig struct {
	Config
	TLSConfig *tls.Config
	NetDialer *net.Dialer
}

// Dial dials a connection.
func (d DialConfig) Dial(ctx context.Context) (conn net.Conn, err error) {
	if conn, err = d.NetDialer.DialContext(ctx, d.Network, d.Address); err != nil {
		return
	}
	if d.TLSConfig != nil {
		conn = tls.Client(conn, d.TLSConfig)
	}
	return
}

// ParseDialURI parses a dial URI into a DialConfig.
// The URI must be in the form of:
//
//	[scheme://][address][:port][?option=value[&option=value]]
//
// See ParseListenURI for supported schemes.
//
// Supported options are:
//   - tls-ca-file: path to CA file of server certificate
//   - tls-cert-file: TLS certificate file path, for client auth
//   - tls-key-file: TLS certificate key file path, for client auth
//   - tls-server-name: TLS server name, overrides connection address
//   - insecure: skip tls verification
func ParseDialURI(s string, overrides ...map[string]string) (cfg DialConfig, err error) {
	cfg.NetDialer = &net.Dialer{}

	if cfg.Config, err = ParseURI(s, overrides...); err != nil {
		return
	}

	if cfg.Secure {
		cfg.TLSConfig = &tls.Config{}

		for k, v := range cfg.Options {
			switch k {
			case OptionTLSServerName:
				cfg.TLSConfig.ServerName = v
			case OptionTLSCAFile:
				if cfg.TLSConfig.RootCAs, err = certPoolWithFile(v, true); err != nil {
					return
				}
			case OptionInsecure:
				cfg.TLSConfig.InsecureSkipVerify, _ = strconv.ParseBool(v)
			}
		}

		if cfg.TLSConfig.ServerName == "" && !cfg.TLSConfig.InsecureSkipVerify {
			if host, _, _ := net.SplitHostPort(cfg.Address); host != "" {
				cfg.TLSConfig.ServerName = host
			}
		}

		// client auth
		if crtFile, keyFile := cfg.Options[OptionTLSCertFile], cfg.Options[OptionTLSKeyFile]; crtFile != "" && keyFile != "" {
			var crt tls.Certificate
			if crt, err = tls.LoadX509KeyPair(crtFile, keyFile); err != nil {
				return
			}
			cfg.TLSConfig.Certificates = append(cfg.TLSConfig.Certificates, crt)
		}
	}

	return
}
