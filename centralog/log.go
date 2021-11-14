package centralog

import (
	"fmt"
	"net"
	"net/url"
	"os"
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

type logSink struct {
	ws zapcore.WriteSyncer
	f  *os.File
	sw *serverWriter
}

func (l *logSink) Write(p []byte) (n int, e error) {
	return l.ws.Write(p)
}

func (l *logSink) Sync() error {
	return l.ws.Sync()
}

func (l *logSink) Close() error {
	err := l.f.Close()
	if err != nil {
		return err
	}
	err = l.sw.Close()
	if err != nil {
		return err
	}
	return nil
}

func newLogSink(client net.Conn, f *os.File) zap.Sink {
	sw := newServerWriter(client)
	return &logSink{
		ws: zapcore.NewMultiWriteSyncer(f, sw), // !important:DO NOT change this order of params
		f:  f,
		sw: sw,
	}
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
		return newLogSink(client, os.Stdout), nil
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
