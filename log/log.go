package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	native "log"
	"strings"
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
	ShowMasterKey = "ShowMaster"
	GameEventKey  = "GameEvent"
)

type iStorage interface {
	Add(area, key string, val string) error
	NextID() int
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
		entry.Debug(args...)
	case LVL_INFO:
		logger.Info(args...)
	case LVL_WARN:
		logger.Warn(args...)
	case LVL_ERROR:
		logger.Error(args...)
	case LVL_PANIC:
		logger.Panic(args...)
	case LVL_FATAL:
		logger.Fatal(args...)
	}
}

func Log(lvl int, args ...interface{}) {
	if logger == nil {
		if lvl <= int(LoggerLevel) {
			native.Println(args...)
		}
		return
	}

	entry := logrus.NewEntry(logger)
	sfmu.RLock()
	entry = entry.WithFields(stateFields)
	sfmu.RUnlock()
	log(lvl, entry, args...)
}

func LogGame(key string, args ...interface{}) {
	if logger == nil {
		native.Println("key = ", key, args)
		return
	} else {
		entry := logrus.NewEntry(logger)
		sfmu.RLock()
		entry = entry.WithFields(stateFields)
		sfmu.RUnlock()
		entry = entry.WithField(ShowMasterKey, true)
		entry = entry.WithField(GameEventKey, key)
		entry.Info(args...)
	}
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
	t1 := time.Now().UnixNano()
	return func() {
		d := (time.Now().UnixNano() - t1) / 1000
		Log(LVL_DEBUG, "Func: ", name, "ended time:", d, "McS")
	}
}

func SetStorage(storage iStorage) {
	logStorage = storage
}

func GetLogStateFieldsStr() string {
	sfmu.Lock()
	defer sfmu.Unlock()
	return fmt.Sprint(stateFields)
}

func SaveToStorage(eventKey string, args string, stateFields string) {
	forbidden := []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", "+"}

	area := eventKey
	nID := logStorage.NextID()
	tStr := time.Now().Format(time.Stamp)
	tStr = strings.Replace(tStr, ":", ".", -1)
	safeArgs := args
	for _, sym := range forbidden {
		safeArgs = strings.Replace(safeArgs, sym, ".", -1)
	}

	key := fmt.Sprintf("%v at %v (%v)", safeArgs, tStr, nID)

	err := logStorage.Add(area, key, stateFields+"\n"+args)
	if err != nil {
		Log(LVL_ERROR, "save to storage error", err)
	}
}
