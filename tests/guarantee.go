package main

import (
	"fmt"
	"github.com/Shnifer/magellan/wrnt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var mu sync.Mutex
var s *wrnt.SendMany

var names = [3]string{"alice", "bob", "candy"}
var r [3]*wrnt.Recv

var dataIn [3]chan wrnt.Storage
var dataOut [3]chan wrnt.Storage
var confIn [3]chan int
var confOut [3]chan int

func NewGenerator() func() string {
	var i int
	return func() string {
		i++
		return strconv.Itoa(i)
	}
}

func Generator() {
	gen := NewGenerator()

	for {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second / 10)
		x := gen()
		log.Println("sent " + x)
		//fmt.Println("Generator lock")
		mu.Lock()
		s.AddItems(x)
		//fmt.Println("Generator unlock")
		mu.Unlock()
	}
}

func Sender(n int) {

	for {
		time.Sleep(time.Second / 2)
		//fmt.Println("Sender lock")
		mu.Lock()
		msg, err := s.Pack(names[n])
		//fmt.Println("Sender unlock")
		mu.Unlock()
		if err == nil {
			dataIn[n] <- msg
		}
	}
}

func Recv(n int) {
	t := time.Tick(time.Second / 2)
	for {
		select {
		case data := <-dataOut[n]:
			//fmt.Println("Recv data lock")
			mu.Lock()
			msg := r[n].Unpack(data)
			//fmt.Println("Recv data unlock")
			mu.Unlock()
			for _, v := range msg {
				log.Println(names[n], "receive "+v)
			}
		case <-t:
			//fmt.Println("Recv conf lock")
			mu.Lock()
			conf := r[n].LastRecv()
			//fmt.Println("Recv conf unlock")
			mu.Unlock()
			confIn[n] <- conf

		}
	}
}

func BadMedia(n int) {
	for {
		select {
		case data := <-dataIn[n]:
			if rand.Intn(3) > 0 {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second / 10)
				log.Println(names[n], data.Items)
				dataOut[n] <- data
			} else {
				log.Println(names[n], "data lost")
			}
		case N := <-confIn[n]:
			if rand.Intn(3) > 0 {
				time.Sleep(time.Duration(rand.Intn(10)) * time.Second / 10)
				log.Println(names[n], "confirmed", N)
				confOut[n] <- N
			} else {
				log.Println(names[n], "conf lost")
			}
		}
	}
}

func Confer(n int) {
	for {
		N := <-confOut[n]
		log.Println(names[n], "conf get", N)
		//fmt.Println("Confer lock")
		mu.Lock()
		s.Confirm(names[n], N)
		//fmt.Println("Confer unlock")
		mu.Unlock()
	}
}

//TODO: why blocks???
const treads = 1

func main() {
	rand.Seed(time.Now().Unix())

	s = wrnt.NewSendMany(names[:treads])

	for i := 0; i < treads; i++ {
		dataIn[i] = make(chan wrnt.Storage, 1)
		dataOut[i] = make(chan wrnt.Storage, 1)
		confIn[i] = make(chan int, 1)
		confOut[i] = make(chan int, 1)
		r[i] = wrnt.NewRecv()
	}

	go Generator()

	for i := 0; i < treads; i++ {
		go Sender(i)
		go Recv(i)
		go BadMedia(i)
		go Confer(i)
	}

	for {
		time.Sleep(3 * time.Second)
		s.Reset()
		log.Println("RESET!!!")
	}

	var s string
	fmt.Scanln(&s)
}
