package centralog

import (
	"context"

	"go.uber.org/zap"
)

const (
	IDKey string = "centralog_id"
)

type CtxProvider interface {
	context.Context
	Set(key string, value interface{})
	Get(key string) (value interface{}, exists bool)
	Iterate(func(key string, value interface{}))
}

type ChainLogger struct {
	logFunc
	msg    string
	fields []zap.Field
}

func (cl *ChainLogger) Log() {
	cl.logFunc(cl.msg, cl.fields...)
}
func (cl *ChainLogger) Any(key string, val interface{}) *ChainLogger {
	cl.fields = append(cl.fields, zap.Any(key, val))
	return cl
}
func (cl *ChainLogger) CtxID(ctx context.Context) *ChainLogger {
	return cl.Any("id", ctx.Value(IDKey))
}

type logFunc func(msg string, fields ...zap.Field)

func Debug(msg string) *ChainLogger {
	return &ChainLogger{
		logFunc: defaultLogger.Debug,
		msg:     msg,
	}
}
func Info(msg string) *ChainLogger {
	return &ChainLogger{
		logFunc: defaultLogger.Info,
		msg:     msg,
	}
}
func Warn(msg string) *ChainLogger {
	return &ChainLogger{
		logFunc: defaultLogger.Warn,
		msg:     msg,
	}
}
func Error(msg string) *ChainLogger {
	return &ChainLogger{
		logFunc: defaultLogger.Error,
		msg:     msg,
	}
}
func Sync() {
	defaultLogger.Sync()
}

func CtxDebug(ctx CtxProvider, msg string) *ChainLogger {
	c := &ChainLogger{
		logFunc: defaultLogger.Debug,
		msg:     msg,
	}
	ctx.Iterate(func(key string, value interface{}) {
		c.fields = append(c.fields, zap.Any(key, value))
	})
	return c
}

func CtxInfo(ctx CtxProvider, msg string) *ChainLogger {
	c := &ChainLogger{
		logFunc: defaultLogger.Info,
		msg:     msg,
	}
	ctx.Iterate(func(key string, value interface{}) {
		c.fields = append(c.fields, zap.Any(key, value))
	})
	return c
}

func CtxWarn(ctx CtxProvider, msg string) *ChainLogger {
	c := &ChainLogger{
		logFunc: defaultLogger.Warn,
		msg:     msg,
	}
	ctx.Iterate(func(key string, value interface{}) {
		c.fields = append(c.fields, zap.Any(key, value))
	})
	return c
}

func CtxError(ctx CtxProvider, msg string) *ChainLogger {
	c := &ChainLogger{
		logFunc: defaultLogger.Error,
		msg:     msg,
	}
	ctx.Iterate(func(key string, value interface{}) {
		c.fields = append(c.fields, zap.Any(key, value))
	})
	return c
}
