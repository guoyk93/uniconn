package uniconn

import (
	"context"
	"github.com/stretchr/testify/require"
	"net"
	"path/filepath"
	"testing"
	"time"
)

func TestListenConfig_Listen(t *testing.T) {
	t.Run("listen-tcp", func(t *testing.T) {
		cfg, err := ParseListenURI("tcp://127.0.0.1:29990")
		require.NoError(t, err)
		lis, err := cfg.Listen(context.Background())
		require.NoError(t, err)
		defer lis.Close()

		go func() {
			if conn, err := lis.Accept(); err == nil {
				_, _ = conn.Write([]byte("OK"))
			}
		}()

		conn, err := net.Dial("tcp", "127.0.0.1:29990")
		require.NoError(t, err)

		buf := make([]byte, 8, 8)
		n, err := conn.Read(buf)
		require.NoError(t, err)
		require.Equal(t, "OK", string(buf[:n]))
		conn.Close()
	})
}

func TestParseListenURI(t *testing.T) {
	cfg, err := ParseListenURI("tcp://127.0.0.1:8080")
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.NotNil(t, cfg.NetListenConfig)
	require.Nil(t, cfg.TLSConfig)

	cfg, err = ParseListenURI("tcp://127.0.0.1:8080?keep-alive=5m")
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.NotNil(t, cfg.NetListenConfig)
	require.Equal(t, time.Minute*5, cfg.NetListenConfig.KeepAlive)
	require.Nil(t, cfg.TLSConfig)

	cfg, err = ParseListenURI("tcp://127.0.0.1:8080?keep-alive=5m&multipath-tcp=true")
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.NotNil(t, cfg.NetListenConfig)
	require.Equal(t, true, cfg.NetListenConfig.MultipathTCP())
	require.Nil(t, cfg.TLSConfig)

	cfg, err = ParseListenURI("tcp+ssl://127.0.0.1:8080?keep-alive=5m&multipath-tcp=true")
	require.Error(t, err)

	cfg, err = ParseListenURI(
		"tcp+ssl://127.0.0.1:8080?keep-alive=5m&multipath-tcp=true",
		map[string]string{
			OptionCertFile: filepath.Join("testdata", "server.full-crt.pem"),
			OptionKeyFile:  filepath.Join("testdata", "server.key.pem"),
		},
	)
	require.NoError(t, err)
	require.Equal(t, cfg.Network, "tcp")
	require.Equal(t, cfg.Address, "127.0.0.1:8080")
	require.NotNil(t, cfg.NetListenConfig)
	require.Equal(t, true, cfg.NetListenConfig.MultipathTCP())
	require.NotNil(t, cfg.TLSConfig)
	require.Equal(t, 1, len(cfg.TLSConfig.Certificates))
}
