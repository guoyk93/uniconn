package uniconn

import (
	"context"
	"github.com/stretchr/testify/require"
	"net"
	"net/url"
	"path/filepath"
	"testing"
)

func TestDialConfig_Dial(t *testing.T) {
	type suite struct {
		name       string
		schema     string
		addressL   string
		addressD   string
		serverName bool
	}
	suites := []suite{
		{
			name:     "tcp-tls",
			schema:   "tcp+tls",
			addressL: "127.0.0.1:19990",
			addressD: "localhost:19990",
		},
		{
			name:       "unix-tls",
			schema:     "unix+tls",
			addressL:   "/tmp/uniconn-test-socket",
			addressD:   "/tmp/uniconn-test-socket",
			serverName: true,
		},
	}

	for _, s := range suites {
		t.Run("dial-"+s.name, func(t *testing.T) {
			ql := &url.Values{}
			ql.Set(OptionCertFile, filepath.Join("testdata", "server.full-crt.pem"))
			ql.Set(OptionKeyFile, filepath.Join("testdata", "server.key.pem"))
			ql.Set(OptionClientCAFile, filepath.Join("testdata", "rootca.crt.pem"))

			qd := &url.Values{}
			qd.Set(OptionCAFile, filepath.Join("testdata", "rootca.crt.pem"))
			qd.Set(OptionCertFile, filepath.Join("testdata", "client.full-crt.pem"))
			qd.Set(OptionKeyFile, filepath.Join("testdata", "client.key.pem"))
			if s.serverName {
				qd.Set(OptionServerName, "localhost")
			}

			var (
				listenURI = s.schema + "://" + s.addressL + "?" + ql.Encode()
				dialURI   = s.schema + "://" + s.addressD + "?" + qd.Encode()
			)

			cfgL, err := ParseListenURI(listenURI)
			require.NoError(t, err)
			require.NotNil(t, cfgL.NetListenConfig)
			require.NotNil(t, cfgL.TLSConfig)

			cfgD, err := ParseDialURI(dialURI)
			require.NoError(t, err)
			require.NotNil(t, cfgD.TLSConfig)

			ctx := context.Background()

			lis, err := cfgL.Listen(ctx)
			require.NoError(t, err)
			defer lis.Close()

			go func() {
				for {
					conn, err := lis.Accept()
					if err != nil {
						t.Log(err.Error())
						return
					} else {
						go func(conn net.Conn) {
							defer conn.Close()
							_, err := conn.Write([]byte("OK"))
							if err != nil {
								t.Log(err.Error())
							}
						}(conn)
					}
				}
			}()

			conn, err := cfgD.Dial(ctx)
			require.NoError(t, err)
			defer conn.Close()

			buf := make([]byte, 8, 8)

			n, err := conn.Read(buf)
			require.NoError(t, err)
			require.Equal(t, "OK", string(buf[:n]))
		})
	}
}

func TestParseDialURI(t *testing.T) {
	cfg, err := ParseDialURI("tcp+tls://127.0.0.1:8080", map[string]string{
		"insecure": "true",
	})
	require.NoError(t, err)
	require.NotNil(t, cfg.TLSConfig)
	require.True(t, cfg.TLSConfig.InsecureSkipVerify)
}
