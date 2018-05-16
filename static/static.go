package static

import (
	"bytes"
	"github.com/Shnifer/magellan/commons"
	"github.com/gobuffalo/packr"
	"io"
	"io/ioutil"
	"os"
	"errors"
)

//box path must be just a string to be parsed by
const resBoxPath = "../res/"
const resFilePath = "res/"

var resBox packr.Box

func init() {
	resBox = packr.NewBox(resBoxPath)
}

func Load(pack, filename string) ([]byte, error) {
	fn := pack + "/" + filename
	if _, err := os.Stat(resFilePath + pack + "/" + filename); err == nil {
		commons.Log(commons.LVL_ERROR, "Load", pack, filename, "from external file")
		return ioutil.ReadFile(resFilePath + pack + "/" + filename)
	} else if resBox.Has(fn) {
		commons.Log(commons.LVL_ERROR, "Load", pack, filename, "from embedded")
		return resBox.MustBytes(pack + "/" + filename)
	} else {
		return nil,errors.New("static "+pack+" "+filename+" not found!")
	}
}

func Exist(pack, filename string) bool {
	inBox := resBox.Has(pack + "/" + filename)
	if inBox {
		return true
	}
	if _, err := os.Stat(resFilePath + pack + "/" + filename); err == nil {
		return true
	} else {
		commons.Log(commons.LVL_WARNING, "Check embedded for", pack, filename, "miss")
		return false
	}
}

func Read(pack, filename string) (io.Reader, error) {
	b, err := Load(pack, filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
