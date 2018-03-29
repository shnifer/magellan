package input

import (
	"github.com/hajimehoshi/ebiten"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
)

const (
	keyPressed = "KP"
	keyJustPressed = "KJP"
	joyAxis = "JA"
	joyButtonPressed = "JBP"
	joyButtonJustPressed = "JBJP"
	mouseButtonPressed = "MBP"
	mouseButtonJustPressed = "MBJP"
)

type InputImpl struct{
	InputName string
	InputType string
	Key ebiten.Key
	JoyID int
	JoyAxis int
	JoyButton int
	MouseButton ebiten.MouseButton
}

var inputState map[string]float64

var conf []InputImpl

func Get(inputName string) bool{
	return false
}

func GetF(inputName string) float64{
	return 0.0
}

func Update(){

}

var savedDef bool
func LoadConf(filePath string) {

	if !savedDef {
		DefConf(filePath)
		savedDef = true
	}

	fn := filePath + "input.json"

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &conf)
}

func DefConf(filePath string){
	def:=make([]InputImpl,2)
	def[0] = InputImpl{
		InputName: "forward",
		InputType: keyPressed,
		Key: ebiten.KeyUp,
	}
	def[1] = InputImpl{
		InputName: "turn",
		InputType:joyAxis,
		JoyID: 0,
		JoyAxis: 1,
	}
	exfn := filePath + "example_input.json"
	exbuf, _ := json.Marshal(def)
	identbuf:=bytes.Buffer{}
	json.Indent(&identbuf, exbuf,"","    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	var str string
	add:=func(name, value string) {
		str=str+name+" = "+value+"\n"
	}
	add("keyPressed", keyPressed)
	add("keyJustPressed", keyJustPressed)
	add("joyAxis", joyAxis)
	add("joyButtonPressed", joyButtonPressed)
	add("joyButtonJustPressed", joyButtonJustPressed)
	add("mouseButtonPressed", mouseButtonPressed)
	add("mouseButtonJustPressed", mouseButtonJustPressed)

	add("","")

	for i,name:=range keysNames{
		add(name, strconv.Itoa(i))
	}

	excodesfn:=filePath + "example_input_codes.txt"
	if err:= ioutil.WriteFile(excodesfn, []byte(str), 0); err!=nil{
		log.Println("can't even write ", excodesfn)
	}
}

var keysNames = []string{
"Key0",
"Key1",
"Key2",
"Key3",
"Key4",
"Key5",
"Key6",
"Key7",
"Key8",
"Key9",
"KeyA",
"KeyB",
"KeyC",
"KeyD",
"KeyE",
"KeyF",
"KeyG",
"KeyH",
"KeyI",
"KeyJ",
"KeyK",
"KeyL",
"KeyM",
"KeyN",
"KeyO",
"KeyP",
"KeyQ",
"KeyR",
"KeyS",
"KeyT",
"KeyU",
"KeyV",
"KeyW",
"KeyX",
"KeyY",
"KeyZ",
"KeyAlt",
"KeyApostrophe",
"KeyBackslash",
"KeyBackspace",
"KeyCapsLock",
"KeyComma",
"KeyControl",
"KeyDelete",
"KeyDown",
"KeyEnd",
"KeyEnter",
"KeyEqual",
"KeyEscape",
"KeyF1",
"KeyF2",
"KeyF3",
"KeyF4",
"KeyF5",
"KeyF6",
"KeyF7",
"KeyF8",
"KeyF9",
"KeyF10",
"KeyF11",
"KeyF12",
"KeyGraveAccent",
"KeyHome",
"KeyInsert",
"KeyLeft",
"KeyLeftBracket",
"KeyMinus",
"KeyPageDown",
"KeyPageUp",
"KeyPeriod",
"KeyRight",
"KeyRightBracket",
"KeySemicolon",
"KeyShift",
"KeySlash",
"KeySpace",
"KeyTab",
"KeyUp",
}
