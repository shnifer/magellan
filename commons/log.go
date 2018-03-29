package commons

import "fmt"

const (
	LVL_DEBUG int = iota
	LVL_WARNING
	LVL_ERROR
	LVL_PANIC
	LVL_NOLOG
)

const (
	LOG_LEVEL = LVL_DEBUG
)

func Log(level int, params ...interface{}) {
	if level >= LOG_LEVEL {
		fmt.Println()
	}
}
