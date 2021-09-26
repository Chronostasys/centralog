package centralog

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
	addr          string
	dbcol         string
)

type LogOptions struct {
	Server     string
	Db         string // optional
	Collection string // optional
}

func InitLoggerWithOpt(conf zap.Config, opts *LogOptions) error {
	err := InitLogger(conf, opts.Server)
	if err != nil {
		return err
	}
	dbcol = opts.Db + ";" + opts.Collection
	selectDbAndCol()
	return nil
}
func selectDbAndCol() {
	Info(dbcol).Log()
	Sync()
}
func InitLogger(conf zap.Config, serverAddr string) error {
	connChan = make(chan net.Conn)
	addr = serverAddr
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}
	conn = client
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
	// Log func in wrapped in centralog
	// so add a default skip to make the caller right
	defaultLogger = defaultLogger.WithOptions(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	return nil
}
