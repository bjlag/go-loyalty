//go:generate mockgen -source ${GOFILE} -package mock -destination mock/logger_mock.go

package logger

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Error(msg string)
	Errorf(template string, args ...interface{})
	Warning(msg string)
	Warningf(template string, args ...interface{})
	Info(msg string)
	Infof(template string, args ...interface{})
	Debug(msg string)
	Debugf(template string, args ...interface{})
}
