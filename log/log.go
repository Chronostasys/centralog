package log

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func InitLogger(conf zap.Config, serverAddr string) error {
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}
	conf.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02T15:04:05Z07:00"))
	}
	zap.RegisterSink("log-server", func(u *url.URL) (zap.Sink, error) {
		return newServerWriter(client), nil
	})
	// build a valid custom path
	customPath := fmt.Sprintf("%s:whatever", "log-server")
	conf.OutputPaths = []string{customPath}
	defaultLogger, err = conf.Build()
	if err != nil {
		return err
	}
	return nil
}
