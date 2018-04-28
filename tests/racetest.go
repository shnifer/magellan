package main

import "time"

var i int

func inc(){
	for{
		time.Sleep(time.Microsecond)
		i++
	}
}

func dec(){
	for{
		time.Sleep(time.Microsecond)
		i--
	}
}

func main(){
	go inc()
	go dec()

	time.Sleep(2*time.Second)
}