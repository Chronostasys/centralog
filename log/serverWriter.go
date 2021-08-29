package log

import (
	"bytes"
	"io"
	"net"
)

type serverWriter struct {
	b    *bytes.Buffer
	conn net.Conn
}

func (sw *serverWriter) Close() error {
	return sw.conn.Close()
}
func (sw *serverWriter) Sync() error {
	bs, err := io.ReadAll(sw.b)
	if err != nil {
		return err
	}
	_, err = sw.conn.Write(bs)
	if err != nil {
		return err
	}
	return nil
}

func (sw *serverWriter) Write(p []byte) (n int, err error) {
	n, err = sw.b.Write(p)
	if err != nil {
		return 0, err
	}
	return
}

func newServerWriter(conn net.Conn) *serverWriter {
	buff := bytes.NewBuffer(make([]byte, 0, 500))
	return &serverWriter{
		b:    buff,
		conn: conn,
	}
}
