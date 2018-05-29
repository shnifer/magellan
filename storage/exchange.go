package storage

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

const exchangerPath = "/exchange/"

type neededKeys []string

func (x neededKeys) Len() int { return len(x) }
func (x neededKeys) Less(i, j int) bool {
	if strings.HasPrefix(x[i], glyphDel) &&
		!strings.HasPrefix(x[j], glyphDel) {
		return true
	} else {
		return x[i] < x[j]
	}
}
func (x neededKeys) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type exchanger struct {
	disk *disk

	addrN int
	addrs []string

	sync.RWMutex
	//map[addr] []fullkeys
	needKeys map[string]neededKeys

	client *http.Client
	server *http.Server
}

func RunExchanger(storage *Storage, listenAddr string, addrs []string, periodMs int) {

	if len(addrs) == 0 {
		panic("NewExchanger: addrs must have at least one address")
	}
	if periodMs == 0 {
		panic("NewExchanger: zero periodMs")
	}

	period := time.Duration(periodMs) * time.Millisecond

	client := &http.Client{
		Timeout: period,
	}

	mux := http.NewServeMux()
	server := &http.Server{Addr: listenAddr, Handler: mux}

	res := &exchanger{
		disk:     storage.disk,
		addrs:    addrs,
		client:   client,
		server:   server,
		needKeys: make(map[string]neededKeys, 0),
	}

	mux.Handle(exchangerPath, exchangeHandler(res))

	go func() {
		for {
			time.Sleep(period)
			res.exchange()
		}
	}()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()
}

func (ex *exchanger) exchange() {
	ex.addrN = (ex.addrN + 1) % len(ex.addrs)
	addr := ex.addrs[ex.addrN]

	needed, ok := ex.needKeys[addr]
	if !ok {
		needed = []string{}
		ex.needKeys[addr] = needed
	}

	needKey := ""
	for len(needed) > 0 {
		mbKey := needed[0]
		if !ex.disk.has(mbKey) {
			needKey = mbKey
			break
		}
		needed = needed[1:]
	}

	req := Request{
		INeedFullKey: needKey,
	}
	reqBody := req.Encode()

	respBuf, err := ex.doReq(addr+exchangerPath, reqBody)
	if err != nil {
		//it is ok to have no access
		return
	}

	resp, err := Responce{}.Decode(respBuf)
	if err != nil {
		return
	}

	if needKey != "" {
		err := ex.disk.append(needKey, resp.YourKeyVal)
		if err != nil {
			log.Println("downloaded needKey already exist:", needKey)
		}
	}

	for _, key := range resp.IHaveFullKeys {
		if !ex.disk.has(key) {
			needed = append(needed, key)
		}
	}

	//glyph of delete "!" must be first, to get delete flags first
	sort.Strings(needed)

	ex.needKeys[addr] = needed
}

func exchangeHandler(ex *exchanger) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		reqBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		req, err := Request{}.Decode(reqBuf)
		if err != nil {
			log.Println(err)
			return
		}
		resp := Responce{}
		if req.INeedFullKey != "" {
			val, err := ex.disk.Read(req.INeedFullKey)
			if err != nil {
				log.Println("Can't read fullkey ", req.INeedFullKey)
				return
			}
			resp.YourKeyVal = string(val)
		}

		resp.IHaveFullKeys = ex.disk.getKeys()
		w.Write(resp.Encode())
	}

	return http.HandlerFunc(f)
}

func (ex *exchanger) doReq(addr string, reqBody []byte) (respBody []byte, er error) {
	bodyBuf := bytes.NewBuffer(reqBody)
	req, err := http.NewRequest(http.MethodGet, addr, bodyBuf)
	if err != nil {
		return nil, err
	}

	resp, err := ex.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
