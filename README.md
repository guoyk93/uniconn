# uniconn

A universal connection dial and listen library, support `tcp` and `unix` with TLS

## Usage

```go
package main

import (
	"context"
	"github.com/guoyk93/uniconn"
)

func listen() {
	cfg, err := uniconn.ParseListenURI("tcp+tls://127.0.0.1:8443?cert-file=cert.pem&key-file=key.pem")
	if err != nil {
		panic(err)
	}
	lis, err := cfg.Listen(context.Background())
	if err != nil {
		panic(err)
	}
	lis.Accept()

	//...
}

func dial() {
	cfg, err := uniconn.ParseDialURI("tcp+tls://127.0.0.1:8443?ca-file=ca.pem")
	if err != nil {
		panic(err)
	}
	conn, err := cfg.Dial(context.Background())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// ...
}
```

## Options

**Listen**

- `cert-file`: TLS certificate file path, mandatory if TLS is enabled
- `key-file`: TLS certificate key file path, mandatory if TLS is enabled
- `client-ca-file`: TLS certificate file path to client ca, setting this option will enable client authentication
- `keep-alive`: TCP keep-alive period, in format of "1m", "1h", etc.
- `multipath-tcp`: Enable multipath TCP, in format of "true" or "false"

**Dial**

- `ca-file`: TLS CA file path
- `cert-file`: TLS certificate file path, for client auth
- `key-file`: TLS certificate key file path, for client auth
- `server-name`: TLS server name, overrides connection address
- `insecure`: skip tls verification

## Donation

View https://guoyk.xyz/donation

## Credits

GUO YANKE, MIT License
