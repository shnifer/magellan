package main


import (
	"net"
)

func main(){
	c, err := net.Dial("tcp", "localhost:6666")
	if err!=nil{
		panic(err)
	}
	defer c.Close()

	c.Write([]byte("hello!"))
}
