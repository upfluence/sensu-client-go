package sentry

import (
	"fmt"
	"os"

	"github.com/getsentry/raven-go"
	"github.com/upfluence/goutils/thrift/handler"
)

type ErrorLogger struct {
	client *raven.Client
}

func NewErrorLogger(dsn string) (*ErrorLogger, error) {
	cl, err := raven.NewClient(
		dsn,
		map[string]string{
			"semver_version": handler.Version,
			"git_commit":     handler.GitCommit,
			"git_branch":     handler.GitBranch,
			"git_remote":     handler.GitRemote,
			"unit_name":      os.Getenv("UNIT_NAME"),
		},
	)

	if err != nil {
		return nil, err
	}

	if handler.Version != "v0.0.0" {
		cl.SetRelease(
			fmt.Sprintf("%s-%s", os.Getenv("PROJECT_NAME"), handler.Version),
		)
	}

	if v := os.Getenv("ENV"); v != "" {
		cl.SetEnvironment(v)
	}

	return &ErrorLogger{cl}, nil
}

func (l *ErrorLogger) Capture(err error, opts map[string]interface{}) error {
	var tags = make(map[string]string)

	for k, v := range opts {
		tags[k] = fmt.Sprintf("%+v", v)
	}

	l.client.CaptureError(err, tags)

	return nil
}

func (l *ErrorLogger) Close() {
	l.client.Wait()
	l.client.Close()
}
