package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"sync"
)

var Ini struct {
	MyPort   string
	Addrs    []string
	NodeName string
	MyArea   string
	PeriodMs int
}

func main() {
	buf, err := ioutil.ReadFile("ini.json")
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(buf, &Ini)
	if err != nil {
		log.Panicln(err)
	}

	opts := diskv.Options{
		BasePath: "dat",
	}
	store := storage.New(Ini.NodeName, opts)
	storage.RunExchanger(store, Ini.MyPort, Ini.Addrs, Ini.PeriodMs)

	var datamu sync.Mutex

	fmt.Printf("Node %v online \n", Ini.NodeName)
	data, subscribe := store.SubscribeAndData(Ini.MyArea)
	fmt.Println("subscribed for area", Ini.MyArea)
	for key, val := range data {
		fmt.Printf("stored data: %v == %v", key, val)
	}
	go func() {
		for {
			event := <-subscribe
			t := "ADD"
			if event.Type == storage.Remove {
				t = "DEL"
			}
			fmt.Printf("Subscribe redistred new event: %v for key %v", t, event.Key)

			datamu.Lock()
			if event.Type == storage.Add {
				data[event.Key] = event.Data
			} else {
				delete(data, event.Key)
			}
			datamu.Unlock()
		}
	}()

	s := ""

	for {
		fmt.Scanln(&s)

		switch s {
		//create
		case "1":
			areaName := strconv.Itoa(rand.Intn(5) + 1)
			key := "obj" + strconv.Itoa(store.NextID())
			val := "value for key " + key
			fmt.Printf("create new object for in area %v with key %v \n", areaName, key)
			store.Add(areaName, key, val)

		//delete
		case "2":
			datamu.Lock()
			n := rand.Intn(len(data))
			var delKey storage.ObjectKey
			for key := range data {
				if n == 0 {
					delKey = key
				} else {
					n--
				}

			}
			datamu.Unlock()
			fmt.Printf("delete object with key %v \n", delKey)
			store.Remove(delKey)
		}
	}
}
