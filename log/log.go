package log

import (
	"github.com/sirupsen/logrus"
	native "log"
)

const (
	LVL_DEBUG = int(logrus.DebugLevel)
	LVL_INFO  = int(logrus.InfoLevel)
	LVL_WARN  = int(logrus.WarnLevel)
	LVL_PANIC = int(logrus.PanicLevel)
	LVL_FATAL = int(logrus.FatalLevel)
	LVL_ERROR = int(logrus.ErrorLevel)
)

var stateFields logrus.Fields

func init() {
	stateFields = make(logrus.Fields)
}

func log(lvl int, entry *logrus.Entry, args ...interface{}) {
	switch lvl {
	case LVL_DEBUG:
		entry.Debug(args)
	case LVL_INFO:
		logger.Info(args)
	case LVL_WARN:
		logger.Warn(args)
	case LVL_ERROR:
		logger.Error(args)
	case LVL_PANIC:
		logger.Panic(args)
	case LVL_FATAL:
		logger.Fatal(args)
	}
}

func Log(lvl int, args ...interface{}) {
	if logger == nil {
		native.Println(args)
		return
	}

	entry := logrus.NewEntry(logger)
	entry.WithFields(stateFields)
	log(lvl, entry, args)
}

//TODO: use this in execs
func SetLogFields(keys map[string]string) {
	stateFields = make(logrus.Fields, len(keys))
	for k, v := range keys {
		stateFields[k] = v
	}

	if logger == nil {
		return
	}
	Log(LVL_INFO, "become", keys)
}

func LogGame(key string, args ...interface{}) {
	if logger == nil {
		native.Println("key = ", key, args)
		return
	}

	entry := logrus.NewEntry(logger)
	entry.WithField("Key", key)
	entry.Info(args)
}

func IsLogDebug() bool {
	if logger == nil {
		return false
	}

	return logger.Level == logrus.DebugLevel
}

func SetLogLevel(lvl int) {
	if logger == nil {
		return
	}

	logger.SetLevel(logrus.Level(lvl))
}

func LogFunc(name string) func() {
	if logger == nil {
		return func() {}
	}

	if logger.Level != logrus.DebugLevel {
		return func() {}
	}

	Log(LVL_DEBUG, "Func: ", name, "started")
	return func() {
		Log(LVL_DEBUG, "Func: ", name, "ended")
	}
}
