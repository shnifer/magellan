package main

import (
	"net/http"
	"log"
	"io/ioutil"
)

func testHandler (w http.ResponseWriter, r *http.Request){
	str:="just русская responce"
	buf:=[]byte(str)
	w.Write(buf)
	if r.Method==http.MethodPost{
		l:=r.ContentLength
		log.Println("req ContentLength",l)

		buf,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			log.Println("POST REQ read err:",err)
		} else {
			log.Println("loaded POST REQ", len(buf))
		}
		str:=string(buf)
		log.Println(len(str), str)
	}
}

func main(){
	http.HandleFunc("/test/", testHandler)
	http.ListenAndServe("localhost:8000",nil)
}