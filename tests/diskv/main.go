package main

import (
	"github.com/peterbourgon/diskv"
	"log"
)

func main() {
	transform := func(s string) []string { return []string{} }
	d := diskv.New(diskv.Options{
		BasePath:     "dat",
		Transform:    transform,
		CacheSizeMax: 1024 * 1024,
	})

	var err error

	err = d.Write("key2", []byte("new value"))
	if err != nil {
		log.Println("write err: ", err)
	}
	err = d.Write("key1", []byte("value"))
	if err != nil {
		log.Println("write err: ", err)
	}

	for key := range d.Keys(nil) {
		str := d.ReadString(key)
		log.Println(key, str)
	}
}
