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

	log.Println(reqString, len(reqString))
	body:=strings.NewReader(reqString)
	req, err := http.NewRequest(http.MethodPost,"http://localhost:8000/test/" , body)
	if err != nil {
		panic(err)
	}

	resp,err:=client.Do(req)
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()

	l:=req.ContentLength
	buf:=make([]byte,l)
	resp.Body.Read(buf)

	log.Println(string(buf))
}