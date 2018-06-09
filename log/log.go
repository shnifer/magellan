package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	native "log"
	"sync"
	"time"
)

const (
	LVL_DEBUG = int(logrus.DebugLevel)
	LVL_INFO  = int(logrus.InfoLevel)
	LVL_WARN  = int(logrus.WarnLevel)
	LVL_PANIC = int(logrus.PanicLevel)
	LVL_FATAL = int(logrus.FatalLevel)
	LVL_ERROR = int(logrus.ErrorLevel)
)

const (
	GameEventKey = "GameEvent"
)

type iStorage interface {
	Add(area, key string, val string) error
}

var logStorage iStorage

var sfmu sync.RWMutex
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
	sfmu.RLock()
	entry.WithFields(stateFields)
	sfmu.RUnlock()
	log(lvl, entry, args)
}

func SetLogFields(keys map[string]string) {
	defer LogFunc("SetLogFields")()

	sfmu.Lock()
	stateFields = make(logrus.Fields, len(keys))
	for k, v := range keys {
		stateFields[k] = v
	}
	sfmu.Unlock()

	if logger == nil {
		return
	}
	Log(LVL_INFO, "become", keys)
}

func LogGame(key string, args ...interface{}) {
	if logStorage != nil {
		go saveToStorage(key, args)
	}

	if logger == nil {
		native.Println("key = ", key, args)
		return
	} else {
		entry := logrus.NewEntry(logger)
		sfmu.RLock()
		entry.WithFields(stateFields)
		sfmu.RUnlock()
		entry.WithField(GameEventKey, key)
		entry.Info(args)
	}
}

//TODO: check somethere and set dynamically
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
	t1:=time.Now().UnixNano()
	return func() {
		d:=(time.Now().UnixNano()-t1)/1000
		Log(LVL_DEBUG, "Func: ", name, "ended time:",d,"McS")
	}
}

func SetStorage(storage iStorage) {
	logStorage = storage
}

func saveToStorage(eventKey string, args ...interface{}) {
	area := eventKey
	key := fmt.Sprint(args)
	sfmu.RLock()
	val := fmt.Sprint(stateFields)
	sfmu.RUnlock()
	logStorage.Add(area, key, val)
}
