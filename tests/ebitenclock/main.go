package main

import (
	"time"
	"fmt"
	"log"
)

var nowTime int64

func init(){
	nowTime = time.Now().UnixNano()
}

func now() int64{
	return time.Now().UnixNano()
}

func pass(ms int) {
	nowTime+=int64(time.Millisecond)*int64(ms)
}

var frameCh chan bool

func fpsMonitor(){
	buf:=make([]bool,0,60)
	tick:=time.Tick(time.Second)
	var draw, total int
	for {
		select{
		case doDraw:=<-frameCh:
			buf = append(buf,doDraw)
			if doDraw {
				draw++
			}
			total++
		case <-tick:
			str:=""
			for i:=range buf{
				if buf[i]{
					str+="X"
				} else {
					str+="_"
				}
			}
			fmt.Println(str)
			fmt.Println(draw,"/",total)
			draw,total=0,0
			buf = buf[:0]
		}
	}
}

const (
	updateMs = 2
	drawMs = 1
	internalMs = 17
)

var isSlowly bool

func internalDraw() {
	time.Sleep(internalMs*time.Millisecond)
}

func gameF(){
	frameCh<-!isSlowly
	time.Sleep(updateMs*time.Millisecond)
	if !isSlowly{
		time.Sleep(drawMs*time.Millisecond)
	}
}

func mainLoop(){
	for {
		//clock.Update
		count := Update()
		for i := 0; i < count; i++ {
			isSlowly = i < count-1
			gameF()
		}
		internalDraw()
	}
}

func main(){
	r:=(1000-float64(FPS)*float64(updateMs))/float64(drawMs+internalMs)
	log.Println("math max = ",r)
	frameCh = make(chan bool, 1)
	go fpsMonitor()
	mainLoop()
}