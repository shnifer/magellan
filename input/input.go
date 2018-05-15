package input

import (
	"bytes"
	"encoding/json"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"sync"
)

const (
	keyPressed              = "KP"
	keyJustPressed          = "KJP"
	keyJustReleased         = "KJR"
	joyAxis                 = "JA"
	joyButtonPressed        = "JBP"
	joyButtonJustPressed    = "JBJP"
	joyButtonJustReleased   = "JBJR"
	mouseButtonPressed      = "MBP"
	mouseButtonJustPressed  = "MBJP"
	mouseButtonJustReleased = "MBJR"
)

var mu sync.RWMutex
var inputState map[string]float64
var conf []Impl

func init() {
	inputState = make(map[string]float64)
}

type Impl struct {
	InputName   string
	InputType   string
	Key         ebiten.Key
	Val         float64
	JoyID       int
	JoyAxis     int
	MinAxis     float64
	JoyButton   ebiten.GamepadButton
	MouseButton ebiten.MouseButton
}

func (impl Impl) get() float64 {
	switch impl.InputType {
	case keyPressed:
		if ebiten.IsKeyPressed(impl.Key) {
			if impl.Val == 0 {
				return 1
			} else {
				return impl.Val
			}
		}
	case keyJustPressed:
		if inpututil.IsKeyJustPressed(impl.Key) {
			return 1
		}
	case keyJustReleased:
		if inpututil.IsKeyJustReleased(impl.Key) {
			return 1
		}
	case joyAxis:
		k := impl.Val
		if k == 0 {
			k = 1
		}
		axis := ebiten.GamepadAxis(impl.JoyID, impl.JoyAxis)
		if math.Abs(axis) < impl.MinAxis {
			axis = 0
		}
		return k * axis
	case joyButtonPressed:
		if ebiten.IsGamepadButtonPressed(impl.JoyID, impl.JoyButton) {
			return 1
		}
	case joyButtonJustPressed:
		if inpututil.IsGamepadButtonJustPressed(impl.JoyID, impl.JoyButton) {
			return 1
		}
	case joyButtonJustReleased:
		if inpututil.IsGamepadButtonJustReleased(impl.JoyID, impl.JoyButton) {
			return 1
		}
	case mouseButtonPressed:
		if ebiten.IsMouseButtonPressed(impl.MouseButton) {
			return 1
		}
	case mouseButtonJustPressed:
		if inpututil.IsMouseButtonJustPressed(impl.MouseButton) {
			return 1
		}
	case mouseButtonJustReleased:
		if inpututil.IsMouseButtonJustReleased(impl.MouseButton) {
			return 1
		}
	default:
		panic(errors.New("Unknown input type " + impl.InputType))
	}
	return 0
}

func Get(inputName string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return inputState[inputName] > 0
}

func GetF(inputName string) float64 {
	mu.RLock()
	defer mu.RUnlock()
	return inputState[inputName]
}

func Update() {
	mu.Lock()
	defer mu.Unlock()
	for key := range inputState {
		inputState[key] = 0.0
	}
	for _, impl := range conf {
		in := impl.InputName
		cur := inputState[in]
		val := impl.get()
		if math.Abs(val) > math.Abs(cur) {
			inputState[impl.InputName] = val
		}
	}
}

var savedDef bool

func init() {
	//disable example saving
	savedDef = true
}

func LoadConf(fileName string) {
	mu.Lock()
	defer mu.Unlock()

	if !savedDef {
		defConf(fileName)
		savedDef = true
	}

	fn := fileName

	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("cant read ", fn, "using default")
		return
	}
	json.Unmarshal(buf, &conf)
}

func defConf(filePath string) {
	def := make([]Impl, 2)
	def[0] = Impl{
		InputName: "forward",
		InputType: keyPressed,
		Key:       ebiten.KeyUp,
	}
	def[1] = Impl{
		InputName: "turn",
		InputType: joyAxis,
		JoyID:     0,
		JoyAxis:   1,
	}
	exfn := filePath + "example_input.json"
	exbuf, _ := json.Marshal(def)
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, exbuf, "", "    ")
	if err := ioutil.WriteFile(exfn, identbuf.Bytes(), 0); err != nil {
		log.Println("can't even write ", exfn)
	}

	var str string
	add := func(name, value string) {
		str = str + name + " = " + value + "\n"
	}
	add("keyPressed", keyPressed)
	add("keyJustPressed", keyJustPressed)
	add("keyJustReleased", keyJustReleased)
	add("joyAxis", joyAxis)
	add("joyButtonPressed", joyButtonPressed)
	add("joyButtonJustPressed", joyButtonJustPressed)
	add("joyButtonJustReleased", joyButtonJustReleased)
	add("mouseButtonPressed", mouseButtonPressed)
	add("mouseButtonJustPressed", mouseButtonJustPressed)
	add("mouseButtonJustReleased", mouseButtonJustReleased)

	add("", "")

	for i, name := range keysNames {
		add(name, strconv.Itoa(i))
	}

	excodesfn := filePath + "example_input_codes.txt"
	if err := ioutil.WriteFile(excodesfn, []byte(str), 0); err != nil {
		log.Println("can't even write ", excodesfn)
	}
}

//just for examples
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
	"KeyKP0",
	"KeyKP1",
	"KeyKP2",
	"KeyKP3",
	"KeyKP4",
	"KeyKP5",
	"KeyKP6",
	"KeyKP7",
	"KeyKP8",
	"KeyKP9",
	"KeyKPAdd",
	"KeyKPDecimal",
	"KeyKPDivide",
	"KeyKPEnter",
	"KeyKPEqual",
	"KeyKPMultiply",
	"KeyKPSubtract",
	"KeyLeft",
	"KeyLeftBracket",
	"KeyMenu",
	"KeyMinus",
	"KeyNumLock",
	"KeyPageDown",
	"KeyPageUp",
	"KeyPause",
	"KeyPeriod",
	"KeyPrintScreen",
	"KeyRight",
	"KeyRightBracket",
	"KeyScrollLock",
	"KeySemicolon",
	"KeyShift",
	"KeySlash",
	"KeySpace",
	"KeyTab",
	"KeyUp",
}
