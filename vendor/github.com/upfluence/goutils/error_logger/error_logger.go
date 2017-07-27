package error_logger

import (
	"os"

	"github.com/upfluence/goutils/error_logger/noop"
	"github.com/upfluence/goutils/error_logger/opbeat"
	"github.com/upfluence/goutils/error_logger/sentry"
)

var DefaultErrorLogger ErrorLogger

type ErrorLogger interface {
	Capture(error, map[string]interface{}) error
	Close()
}

func init() {
	if v := os.Getenv("SENTRY_DSN"); v != "" {
		l, err := sentry.NewErrorLogger(v)

		if err != nil {
			DefaultErrorLogger = noop.NewErrorLogger()
		} else {
			DefaultErrorLogger = l
		}
	} else if v := os.Getenv("OPBEAT_APP_ID"); v != "" {
		DefaultErrorLogger = opbeat.NewErrorLogger()
	} else {
		DefaultErrorLogger = noop.NewErrorLogger()
	}

	if e := recover(); e != nil {
		if err, ok := e.(error); ok {
			DefaultErrorLogger.Capture(err, nil)
			DefaultErrorLogger.Close()
			panic(err.Error())
		}
	}
}

func Capture(err error, opts map[string]interface{}) error {
	return DefaultErrorLogger.Capture(err, opts)
}

func Close() {
	DefaultErrorLogger.Close()
}
