package alice

import (
	"net/http"
	"time"
	"encoding/base64"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"io/ioutil"
)

var opts Opts
var client *http.Client

type Event struct{
	EvType string `json:"eventType"`
	Data [7]int `json:"data"`
}

type Events []Event

type Opts struct{
	Addr string
	Path string
	Login string
	Password string

	logPass64 string
	url string
}

func InitAlice(initOpts Opts) {
	opts = initOpts

	b:=&bytes.Buffer{}
	enc:=base64.NewEncoder(base64.StdEncoding,b)
	enc.Write([]byte(opts.Login+":"+opts.Password))
	enc.Close()

	opts.logPass64 = b.String()
	opts.url = opts.Addr+"/"+opts.Path+"/"

	client = &http.Client{
		Timeout: time.Second,
	}
}

type reqData struct{
	Events `json:"events"`
}

func DoReq(location string, events Events) error{

	var data reqData
	data.Events = events
	body,err :=json.Marshal(data)
	if err!=nil{
		return err
	}
	log.Println(string(body))
	bodyBuf :=bytes.NewBuffer(body)
	req,err:=http.NewRequest(http.MethodPost, opts.url+location, bodyBuf)
	if err!=nil{
		return err
	}
	req.Header.Add("Accept","application/json")
	req.Header.Add("Content-Type","application/json")
	req.Header.Add("Authorization", "Basic "+opts.logPass64)
	resp, err := client.Do(req)
	if err!=nil{
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode!=200{
		return errors.New(resp.Status)
	}
	r,_:=ioutil.ReadAll(resp.Body)
	log.Println(string(r))
	return nil
}