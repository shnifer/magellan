package commons

import (
	"bytes"
	"errors"
	. "github.com/shnifer/magellan/log"
	"io/ioutil"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Second,
	}
}

func DoReq(Method string, addr string, body []byte) (respBody []byte, err error) {
	LogGame("tryReqs", false, Method, addr, string(body))

	bodyBuf := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, addr, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
