package main

import (
	"encoding/json"
	"fmt"
	. "github.com/Shnifer/magellan/log"
	"bytes"
	"net/http"
	"time"
	"errors"
)

const TriesCount = 10

type ReportHy struct {
	Planet   string   `json:"entity_id"`
	Owner    string   `json:"company"`
	Minerals []string `json:"resources"`
}

func reportHyMine(planet string, corp string, mins []int) {
	minerals := make([]string, 0)
	for _, n := range mins {
		minerals = append(minerals, fmt.Sprintf("m%v", n))
	}
	Report := ReportHy{
		Planet:   planet,
		Owner:    corp,
		Minerals: minerals,
	}
	dat, err := json.Marshal(Report)
	if err != nil {
		Log(LVL_ERROR, err)
		return
	}
	pause:=time.Second
	for i:=0;i<TriesCount;i++{
		err:=doReq(dat, DEFVAL.ReportHyMineAddr)
		if err==nil{
			break
		}
		Log(LVL_WARN, "Report hy mine request error: ",err)
		time.Sleep(pause)
		pause*=2
	}
}

func doReq(body []byte, addr string) error {
	client := &http.Client{
		Timeout: time.Second,
	}
	bodyBuf := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, addr, bodyBuf)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return nil
}
