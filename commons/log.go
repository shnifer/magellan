package commons

import "fmt"

const (
	LVL_DEBUG int = iota
	LVL_WARNING
	LVL_ERROR
	//LVL_PANIC
	//LVL_NOLOG
)

const (
	LOG_LEVEL = LVL_ERROR
)

func Log(level int, params ...interface{}) {
	if level >= LOG_LEVEL {
		fmt.Println(params...)
	}
}

//use as
// func  myFunc{
//     defer LogFunc("myFunc")()
// ...
//}
func LogFunc(name string) func() {
	Log(LVL_WARNING, "Func: ", name, "start")
	return func() {
		Log(LVL_WARNING, "Func: ", name, "ended")
	}
}
