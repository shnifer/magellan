package main

import (
	"image"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

//x,y in Galaxy units
func CreateRandomPointGenerator(dens func(int, int) byte) func() image.Point {
	//init
	stop := timer("init RPG")
	defer stop()

	const NUM_ROUTINES = 4

	type calcRes struct {
		x   int
		sum int
	}
	calcXrow := func(x int, r chan<- calcRes, w <-chan bool) {
		sum := 0
		for y := -GalaxyRadius; y < GalaxyRadius; y++ {
			sum += int(dens(x, y))
		}
		r <- calcRes{x, sum}
		<-w
	}

	sums := make([]int, GalaxyRadius*2)
	rCh := make(chan calcRes)
	wCh := make(chan bool, NUM_ROUTINES)

	go func() {
		for x := -GalaxyRadius; x < GalaxyRadius; x++ {
			wCh <- true
			go calcXrow(x, rCh, wCh)
		}
	}()

	for i := -GalaxyRadius; i < GalaxyRadius; i++ {
		CR := <-rCh
		sums[CR.x+GalaxyRadius] = CR.sum
	}

	sum := 0
	for i, v := range sums {
		sum += v
		sums[i] = sum
	}

	RPG := func() image.Point {
		N := rand.Intn(sum)
		X, Y := 0, 0
		for x, s := range sums {
			if s > N {
				if x > 0 {
					N -= sums[x-1]
				}
				X = x - GalaxyRadius
				break
			}
		}
		for y := -GalaxyRadius; y < GalaxyRadius; y++ {
			N -= int(dens(X, y))
			if N < 0 {
				Y = y
				break
			}
		}
		return image.Pt(X, Y)
	}
	return RPG
}
