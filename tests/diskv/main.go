package main

import (
	"github.com/peterbourgon/diskv"
	"log"
	"strings"
)

func AdvancedTransformExample(key string) *diskv.PathKey {
	path := strings.Split(key, "/")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last] + ".txt",
	}
}

// If you provide an AdvancedTransform, you must also provide its
// inverse:

func InverseTransformExample(pathKey *diskv.PathKey) (key string) {
	if len(pathKey.FileName)<4{
		return ""
	}
	txt := pathKey.FileName[len(pathKey.FileName)-4:]
	if txt != ".txt" {
		panic("Invalid file found in storage folder!")
	}
	return strings.Join(pathKey.Path, "/") + pathKey.FileName[:len(pathKey.FileName)-4]
}

func main(){
	d:=diskv.New(diskv.Options{
		BasePath:"dat",
		AdvancedTransform:AdvancedTransformExample,
		InverseTransform:InverseTransformExample,
		CacheSizeMax:1024*1024,
	})

	var err error

	err=d.Write("key2",[]byte("new value"))
	if err!=nil{
		log.Println("write err: ",err)
	}
	err=d.Write("key1",[]byte("value"))
	if err!=nil{
		log.Println("write err: ",err)
	}

	for key:=range d.Keys(nil) {
		str:=d.ReadString(key)
		log.Println(key, str)
	}
}
