package uniconn

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"path/filepath"
	"testing"
)

func TestParseURI(t *testing.T) {
	cfg, err := ParseURI("tcp://127.0.0.1:8080")
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.False(t, cfg.Secure)
	require.Equal(t, map[string]string{}, cfg.Options)

	cfg, err = ParseURI("tcp://127.0.0.1:8080?keep-alive=5m")
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.False(t, cfg.Secure)
	require.Equal(t, map[string]string{"keep-alive": "5m"}, cfg.Options)

	cfg, err = ParseURI("tcp://127.0.0.1:8080?keep-alive=5m&multipath-tcp=true")
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.False(t, cfg.Secure)
	require.Equal(t, map[string]string{"keep-alive": "5m", "multipath-tcp": "true"}, cfg.Options)

	cfg, err = ParseURI(
		"tcp+ssl://127.0.0.1:8080?keep-alive=5m&multipath-tcp=true&cert-file="+url.QueryEscape(filepath.Join("testdata", "ssss")),
		map[string]string{
			OptionCertFile: filepath.Join("testdata", "example-com.crt"),
			OptionKeyFile:  filepath.Join("testdata", "example-com.key"),
		},
	)
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.True(t, cfg.Secure)
	require.Equal(t, map[string]string{"keep-alive": "5m", "multipath-tcp": "true",
		OptionCertFile: filepath.Join("testdata", "example-com.crt"),
		OptionKeyFile:  filepath.Join("testdata", "example-com.key"),
	}, cfg.Options)
}
