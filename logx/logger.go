package logx

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// SetLevel returns logrus.SetLevel.
func SetLevel(level logrus.Level) { logger.SetLevel(level) }

// WithField returns logrus.WithField.
func WithField(key string, value interface{}) *logrus.Entry { return logger.WithField(key, value) }

// WithFields returns logrus.WithFields.
func WithFields(fields logrus.Fields) *logrus.Entry { return logger.WithFields(fields) }

// WithError returns logrus.WithError.
func WithError(err error) *logrus.Entry { return logger.WithError(err) }

// Debugf returns logrus.Debugf.
func Debugf(format string, args ...interface{}) { logger.Debugf(format, args...) }

// Infof returns logrus.Infof.
func Infof(format string, args ...interface{}) { logger.Infof(format, args...) }

// Printf returns logrus.Printf.
func Printf(format string, args ...interface{}) { logger.Printf(format, args...) }

// Warnf returns logrus.Warnf.
func Warnf(format string, args ...interface{}) { logger.Warnf(format, args...) }

// Warningf returns logrus.Warningf.
func Warningf(format string, args ...interface{}) { logger.Warningf(format, args...) }

// Errorf returns logrus.Errorf.
func Errorf(format string, args ...interface{}) { logger.Errorf(format, args...) }

// Fatalf returns logrus.Fatalf.
func Fatalf(format string, args ...interface{}) { logger.Fatalf(format, args...) }

// Panicf returns logrus.Panicf.
func Panicf(format string, args ...interface{}) { logger.Panicf(format, args...) }

// Debug returns logrus.Debug.
func Debug(args ...interface{}) { logger.Debug(args...) }

// Info returns logrus.Info.
func Info(args ...interface{}) { logger.Info(args...) }

// Print returns logrus.Print.
func Print(args ...interface{}) { logger.Print(args...) }

// Warn returns logrus.Warn.
func Warn(args ...interface{}) { logger.Warn(args...) }

// Warning returns logrus.Warning.
func Warning(args ...interface{}) { logger.Warning(args...) }

// Error returns logrus.Error.
func Error(args ...interface{}) { logger.Error(args...) }

// Fatal returns logrus.Fatal.
func Fatal(args ...interface{}) { logger.Fatal(args...) }

// Panic returns logrus.Panic.
func Panic(args ...interface{}) { logger.Panic(args...) }

// Debugln returns logrus.Debugln.
func Debugln(args ...interface{}) { logger.Debugln(args...) }

// Infoln returns logrus.Infoln.
func Infoln(args ...interface{}) { logger.Infoln(args...) }

// Println returns logrus.Println.
func Println(args ...interface{}) { logger.Println(args...) }

// Warnln returns logrus.Warnln.
func Warnln(args ...interface{}) { logger.Warnln(args...) }

// Warningln returns logrus.Warningln.
func Warningln(args ...interface{}) { logger.Warningln(args...) }

// Errorln returns logrus.Errorln.
func Errorln(args ...interface{}) { logger.Errorln(args...) }

// Fatalln returns logrus.Fatalln.
func Fatalln(args ...interface{}) { logger.Fatalln(args...) }

// Panicln returns logrus.Panicln.
func Panicln(args ...interface{}) { logger.Panicln(args...) }

// Tracef returns logrus.Tracef.
func Tracef(format string, args ...interface{}) { logger.Tracef(format, args...) }

// Trace returns logrus.Trace.
func Trace(args ...interface{}) { logger.Trace(args...) }

// Traceln returns logrus.Traceln.
func Traceln(args ...interface{}) { logger.Traceln(args...) }

// Logger defines the available log printers.
type Logger interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	Tracef(format string, args ...interface{})
	Trace(args ...interface{})
	Traceln(args ...interface{})
}
