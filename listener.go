package uniconn

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"strconv"
	"time"
)

// ListenConfig is the configuration for a listener.
type ListenConfig struct {
	Config
	TLSConfig       *tls.Config
	NetListenConfig *net.ListenConfig
}

// Listen listens on the network address and returns a Listener.
func (c ListenConfig) Listen(ctx context.Context) (lis net.Listener, err error) {
	lis, err = c.NetListenConfig.Listen(ctx, c.Network, c.Address)
	if err != nil {
		return
	}
	if c.TLSConfig != nil {
		lis = tls.NewListener(lis, c.TLSConfig)
	}
	return
}

// ParseListenURI parses a listen URI into a ListenConfig.
// The URI must be in the form of:
//
//	[scheme://][address][:port][?option=value[&option=value]]
//
// Supported schemes are:
//   - tcp
//   - tcp4 (tcp)
//   - tcp6 (tcp)
//   - unix
//
// Each schema can add a '+tls' suffix to enable TLS.
//
// Supported options are:
//   - cert-file: TLS certificate file path, mandatory if TLS is enabled
//   - key-file: TLS certificate key file path, mandatory if TLS is enabled
//   - client-ca-file: TLS certificate file path to client ca, setting this option will enable client authentication
//   - keep-alive: TCP keep-alive period, in format of "1m", "1h", etc.
//   - multipath-tcp: Enable multipath TCP, in format of "true" or "false"
func ParseListenURI(s string, overrides ...map[string]string) (cfg ListenConfig, err error) {
	cfg.NetListenConfig = &net.ListenConfig{}

	if cfg.Config, err = ParseURI(s, overrides...); err != nil {
		return
	}

	for k, v := range cfg.Options {
		switch k {
		case OptionKeepAlive:
			var d time.Duration
			if d, err = time.ParseDuration(v); err != nil {
				return
			}
			cfg.NetListenConfig.KeepAlive = d
		case OptionMultipathTCP:
			var b bool
			if b, err = strconv.ParseBool(v); err != nil {
				return
			}
			cfg.NetListenConfig.SetMultipathTCP(b)
		}
	}

	if cfg.Secure {
		cfg.TLSConfig = &tls.Config{}

		if clientCAFile := cfg.Options[OptionClientCAFile]; clientCAFile != "" {
			if cfg.TLSConfig.ClientCAs, err = certPoolWithFile(clientCAFile, false); err != nil {
				return
			}
			cfg.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}

		if crtFile, keyFile := cfg.Options[OptionCertFile], cfg.Options[OptionKeyFile]; crtFile == "" || keyFile == "" {
			err = errors.New("missing tls cert file or key file")
			return
		} else {
			var crt tls.Certificate
			if crt, err = tls.LoadX509KeyPair(crtFile, keyFile); err != nil {
				return
			}
			cfg.TLSConfig.Certificates = append(cfg.TLSConfig.Certificates, crt)
		}
	}
	return
}
