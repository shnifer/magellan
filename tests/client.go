package main

import (
	"net/http"
	"time"
	"log"
	"strconv"
	"strings"
)

func main(){
	client:=http.Client{
		Timeout:time.Second,
	}

	var reqString string
	for i:=1;i<10000;i++{
		reqString=reqString+"my request"+strconv.Itoa(i)
	}

	body:=strings.NewReader(reqString)


	req, err := http.NewRequest(http.MethodPost,"http://192.168.0.1:8000/test/" , body)
	if err != nil {
		switch t:=err.(type){
			default:
				log.Println(t)
				panic(err)
		}
	}

	req.Header.Set("Content-Length", strconv.Itoa(body.Len()))
	resp,err:=client.Do(req)
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()

	l,err:=strconv.Atoi(resp.Header.Get("Content-Length"))
	if err!=nil{
		log.Println("can't read len")
	}
	buf:=make([]byte,l)
	resp.Body.Read(buf)

	log.Println(string(buf))
}