package main

import (
	"net/http"
	"log"
	_"encoding/binary"
	"encoding/binary"
	"fmt"
	"strconv"
)

func testHandler (w http.ResponseWriter, r *http.Request){
	log.Println("METHOD", r.Method)
	log.Println("Header", r.Header)
	log.Println("Body", r.Body)
	str:="just русская responce"
	buf:=[]byte(str)
	l:=uint32(len(buf))
	lbuf:=make([]byte,4)
	binary.BigEndian.PutUint32(lbuf,l)

	w.Header().Set("Content-Length", fmt.Sprint(len(buf)))

	w.Write(buf)

	if r.Method==http.MethodPost{
		l,err:=strconv.Atoi(r.Header.Get("Content-Length"))
		if err!=nil{
			log.Println("CANT read REQ LEN")
		}

		buf:=make([]byte, l)
		r.Body.Read(buf)

		str:=string(buf)
		log.Println(len(str), str)
	}
}

func main(){
	http.HandleFunc("/test/", testHandler)
	http.ListenAndServe("localhost:8000",nil)
}