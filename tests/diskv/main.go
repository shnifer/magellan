package main

import (
	"github.com/peterbourgon/diskv"
	"log"
	"strconv"
	"time"
)

func main() {
	transform := func(s string) []string { return []string{} }
	d := diskv.New(diskv.Options{
		BasePath:     "dat",
		Transform:    transform,
		CacheSizeMax: 1024 * 1024,
	})
	var b []byte

	for i:=1;i<10;i++{
		go func (n int){
			for j:=0;j<10;j++ {
				dat := []byte(strconv.Itoa(n))
				d.Write("thesamekey", dat)
			}
		}(i)
	}

	time.Sleep(time.Second)

	for key := range d.Keys(nil) {
		str := d.ReadString(key)
		log.Println(key, str)
	}
	for key := range d.KeysPrefix("pr",nil) {
		str := d.ReadString(key)
		log.Println(key, str)
	}
}
