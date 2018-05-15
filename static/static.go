package static

import (
	"bytes"
	"fmt"
	"github.com/gobuffalo/packr"
	"io"
	"io/ioutil"
)

//box path must be just a string to be parsed by
const resBoxPath = "../res/"
const resFilePath = "res/"

var resBox packr.Box

func init() {
	resBox = packr.NewBox(resBoxPath)
}

func Load(pack, filename string) ([]byte, error) {
	fmt.Println("Load " + pack + " " + filename)

	fn := pack + "/" + filename
	if resBox.Has(fn) {
		return resBox.MustBytes(pack + "/" + filename)
	} else {
		return ioutil.ReadFile(resFilePath + pack + "/" + filename)
	}
}

func Exist(pack, filename string) bool {
	return resBox.Has(pack + "/" + filename)
}

func Read(pack, filename string) (io.Reader, error) {
	b, err := Load(pack, filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}
