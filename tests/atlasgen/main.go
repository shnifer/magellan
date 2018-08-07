package main

import (
	"fmt"
	"io/ioutil"
)

const maxN = 69
const fstr = "\"HARDPLANET-%v\": {\"FileName\": \"Solid%02d.png\"},\n"

func main(){
	var res string
	for i:=0; i<=maxN;i++{
		res+=fmt.Sprintf(fstr, i,i)
	}

	ioutil.WriteFile("res.txt", []byte(res),0)
}