package ranma

import (
	"encoding/json"
	"fmt"
	. "github.com/shnifer/magellan/log"
	"io/ioutil"
	"strconv"
	"strings"
)

type ReqResp struct {
	Programmed uint16 `json:"programmed"`
	Corrected  uint16 `json:"corrected"`
	Id         int    `json:"timestamp"`
}

func (r *Ranma) recv(sn int) {
	addr := r.addr + strconv.Itoa(sn)
	resp, err := r.client.Get(addr)
	if err != nil {
		Log(LVL_ERROR, "Ranma.recv GET error:", err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(LVL_ERROR, "Ranma.recv read body error", err)
		return
	}
	var reqResp ReqResp
	err = json.Unmarshal(buf, &reqResp)
	if err != nil {
		Log(LVL_ERROR, "Ranma.recv can't unmarshal \"", string(buf), "\" error", err)
	}
	r.corrected[sn].setMsg(reqResp.Corrected, reqResp.Id)
}

func (r *Ranma) send(sn int, x uint16) {
	addr := r.addr + strconv.Itoa(sn)
	req := strings.NewReader(fmt.Sprintf("{\"programmed\":%v}", x))
	resp, err := r.client.Post(addr, "application/json", req)
	if err != nil {
		Log(LVL_ERROR, "Ranma.send POST error:", err)
		return
	}
	resp.Body.Close()
}
