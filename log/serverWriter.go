package log

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
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
	b      *bytes.Buffer
	reConn bool
	m      *sync.Mutex
}

func (sw *serverWriter) Close() error {
	return conn.Close()
}
func (sw *serverWriter) Sync() error {
	bs, err := io.ReadAll(sw.b)
	if err != nil {
		return err
	}
	n, err := getConn().Write(bs)
	fmt.Println(n)
	if err != nil {
		fmt.Println("log server conn interupted, start trying reconnect", err)
		sw.m.Lock()
		if sw.reConn {
			sw.m.Unlock()
			sw.b.Write(bs)
			return err
		}
		sw.reConn = true
		sw.b.Write([]byte(
			fmt.Sprintf(`{"level":"info","ts":"2021-09-26T15:10:19+08:00","caller":"log/log.go:35","msg":"%s"}`, dbcol),
		))
		sw.b.Write(bs)
		sw.m.Unlock()
		go func() {
			for {
				err = reconnectConn()
				if err != nil {
					fmt.Println("reconnect failed, try again after 5 seconds")
					time.Sleep(time.Second * 5)
				} else {
					fmt.Println("reconnected")
					sw.reConn = false
					break
				}

			}
		}()
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
		b:      buff,
		reConn: false,
		m:      &sync.Mutex{},
	}
}
