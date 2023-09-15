package uniconn

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func certPoolWithFile(file string, includeSys bool) (p *x509.CertPool, err error) {
	if includeSys {
		if p, err = x509.SystemCertPool(); err != nil {
			return
		}
	} else {
		p = x509.NewCertPool()
	}

	var buf []byte
	if buf, err = os.ReadFile(file); err != nil {
		return
	}

	for {
		var b *pem.Block
		b, buf = pem.Decode(buf)
		buf = bytes.TrimSpace(buf)
		if b == nil {
			break
		}
		if b.Type == "CERTIFICATE" {
			var crt *x509.Certificate
			if crt, err = x509.ParseCertificate(b.Bytes); err != nil {
				return
			}
			p.AddCert(crt)
		}
	}

	return
}
