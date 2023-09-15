package uniconn

import (
	"context"
	"github.com/quic-go/quic-go"
	"net"
	"time"
)

type quicStreamWrapper struct {
	qc quic.Connection
	qs quic.Stream
	nc net.PacketConn
}

func (q *quicStreamWrapper) Read(b []byte) (int, error) {
	return q.qs.Read(b)
}

func (q *quicStreamWrapper) Write(b []byte) (int, error) {
	return q.qs.Write(b)
}

func (q *quicStreamWrapper) Close() error {
	_ = q.qs.Close()
	return q.qc.CloseWithError(0, "")
}

func (q *quicStreamWrapper) LocalAddr() net.Addr {
	return q.nc.LocalAddr()
}

func (q *quicStreamWrapper) RemoteAddr() net.Addr {
	return q.qc.RemoteAddr()
}

func (q *quicStreamWrapper) SetDeadline(t time.Time) error {
	return q.nc.SetDeadline(t)
}

func (q *quicStreamWrapper) SetReadDeadline(t time.Time) error {
	return q.nc.SetReadDeadline(t)
}

func (q *quicStreamWrapper) SetWriteDeadline(t time.Time) error {
	return q.nc.SetWriteDeadline(t)
}

// WrapQUICServerConn wraps a QUIC connection into a net.Conn by accepting a quic stream.
func WrapQUICServerConn(ctx context.Context, nc net.PacketConn, qc quic.Connection) (wc net.Conn, err error) {
	var qs quic.Stream
	if qs, err = qc.AcceptStream(ctx); err != nil {
		return
	}
	wc = &quicStreamWrapper{
		qc: qc,
		qs: qs,
		nc: nc,
	}
	return
}

// WrapQUICClientConn wraps a QUIC connection into a net.Conn by opening a quic stream.
func WrapQUICClientConn(ctx context.Context, nc net.PacketConn, qc quic.Connection) (wc net.Conn, err error) {
	var qs quic.Stream
	if qs, err = qc.OpenStreamSync(ctx); err != nil {
		return
	}
	wc = &quicStreamWrapper{
		qc: qc,
		qs: qs,
		nc: nc,
	}
	return
}

type quicListenerWrapper struct {
	ctx context.Context
	nc  net.PacketConn
	ql  *quic.Listener
}

func (q *quicListenerWrapper) Accept() (wc net.Conn, err error) {
	var qc quic.Connection
	if qc, err = q.ql.Accept(q.ctx); err != nil {
		return
	}
	if wc, err = WrapQUICServerConn(q.ctx, q.nc, qc); err != nil {
		_ = qc.CloseWithError(0, "")
		return
	}
	return
}

func (q *quicListenerWrapper) Close() error {
	return q.ql.Close()
}

func (q *quicListenerWrapper) Addr() net.Addr {
	return q.ql.Addr()
}

// WrapQUICListener wraps a QUIC listener into a net.Listener.
func WrapQUICListener(ctx context.Context, nc net.PacketConn, ql *quic.Listener) net.Listener {
	return &quicListenerWrapper{
		ctx: ctx,
		nc:  nc,
		ql:  ql,
	}
}
