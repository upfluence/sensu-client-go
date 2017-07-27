package log

import (
	"errors"
	"fmt"

	"github.com/op/go-logging"
	"github.com/upfluence/goutils/error_logger"
)

type errorLoggerBackend struct {
	client error_logger.ErrorLogger
}

func (b *errorLoggerBackend) Log(_ logging.Level, d int, r *logging.Record) error {
	var (
		err    error
		opts   = make(map[string]interface{})
		argIdx int
	)

	if len(r.Args) == 0 {
		return nil
	}

	if err2, ok := r.Args[0].(error); ok {
		argIdx++
		err = err2
	} else {
		err = errors.New(r.Formatted(d + 1))
	}

	for argIdx < len(r.Args)-1 {
		opts[fmt.Sprintf("arg %d", argIdx)] = r.Args[argIdx]
		argIdx++
	}

	return b.client.Capture(err, opts)
}
