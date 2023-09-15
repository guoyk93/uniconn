package uniconn

import (
	"errors"
	"net/url"
	"strings"
)

const (
	OptionTLSServerName   = "tls-server-name"
	OptionTLSCAFile       = "tls-ca-file"
	OptionTLSClientCAFile = "tls-client-ca-file"
	OptionTLSCertFile     = "tls-cert-file"
	OptionTLSKeyFile      = "tls-key-file"
	OptionKeepAlive       = "keep-alive"
	OptionMultipathTCP    = "multipath-tcp"
	OptionInsecure        = "insecure"
)

// Networks is the list of supported network types
var Networks = []string{"tcp", "tcp4", "tcp6", "unix"}

type Config struct {
	Network string
	Address string
	Secure  bool
	Options map[string]string
}

// ParseURI parses generic listen/dial URI
func ParseURI(s string, overrides ...map[string]string) (cfg Config, err error) {
	cfg = Config{
		Options: make(map[string]string),
	}

	splitNetwork := strings.SplitN(s, "://", 2)

	if len(splitNetwork) == 2 {
		cfg.Network = strings.TrimSpace(strings.ToLower(splitNetwork[0]))
		cfg.Address = splitNetwork[1]
	} else if len(splitNetwork) == 1 {
		cfg.Network = "tcp"
		cfg.Address = strings.TrimSpace(splitNetwork[0])
	} else {
		err = errors.New("invalid listen URI")
	}

	if strings.HasSuffix(cfg.Network, "+tls") ||
		strings.HasSuffix(cfg.Network, "+ssl") {
		cfg.Network = cfg.Network[:len(cfg.Network)-4]
		cfg.Secure = true
	}

	for _, network := range Networks {
		if cfg.Network == network {
			goto networkAllowed
		}
	}

	err = errors.New("invalid network type: " + cfg.Network + " (supported: " + strings.Join(Networks, ", "))

	return

networkAllowed:

	if splitOptions := strings.SplitN(cfg.Address, "?", 2); len(splitOptions) == 2 {
		cfg.Address = splitOptions[0]
		var values url.Values
		if values, err = url.ParseQuery(strings.TrimSpace(splitOptions[1])); err != nil {
			return
		}
		for k := range values {
			vs := values[k]
			if len(vs) > 0 {
				cfg.Options[k] = vs[len(vs)-1]
			}
		}
	}

	for _, override := range overrides {
		for k, v := range override {
			if v != "" {
				cfg.Options[k] = v
			}
		}
	}
	return
}
