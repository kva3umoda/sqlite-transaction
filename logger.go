package sqlite_transaction

import (
	"github.com/doug-martin/goqu/v9"
)

//

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

var _ goqu.Logger = (*goquLogger)(nil)

type goquLogger struct {
	logger Logger
}

func newGoquLogger(logger Logger) *goquLogger {
	return &goquLogger{
		logger: logger,
	}
}

func (l *goquLogger) Printf(format string, v ...interface{}) {
	l.logger.Debug(format, v...)
}
