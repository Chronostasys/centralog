package log

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

var (
	conn     net.Conn
	connChan chan net.Conn
)

func reconnectConn() error {
	client, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	connChan <- client
	return nil
}
func getConn() net.Conn {
	select {
	case conn = <-connChan:
	default:
	}
	return conn
}

type serverWriter struct {
	b *bytes.Buffer
}

func (sw *serverWriter) Close() error {
	return conn.Close()
}
func (sw *serverWriter) Sync() error {
	bs, err := io.ReadAll(sw.b)
	if err != nil {
		return err
	}
	_, err = getConn().Write(bs)
	if err != nil {
		sw.b.Write(bs)
		if strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
			fmt.Println("log server conn closed by remote, start trying reconnect")
			go func() {

				for {
					err = reconnectConn()
					if err != nil {
						fmt.Println("reconnect failed, try again after 5 seconds")
						time.Sleep(time.Second * 5)
					} else {
						fmt.Println("reconnected")
						break
					}

				}
			}()
		}
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
		b: buff,
	}
}
