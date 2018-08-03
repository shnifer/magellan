package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
)

type EmiChance struct {
	First, More int
	Far, Close  float64
}

//map[emiN]chanceBasePart
type EmiDistib map[int]int

func (d EmiDistib) gen() int {
	sum := 0

	for _, v := range d {
		if v > 0 {
			sum += v
		}
	}
	if sum == 0 {
		panic("zero distrib base")
	}
	r := rand.Intn(sum)

	i := 0
	for code, base := range d {
		i += base
		if r < i {
			return code
		}
	}

	return 0
}

type Options struct {
	Dev    float64
	Chance struct {
		Star, Asteroid, Hard, Gas, Warp EmiChance
	}
	Star, Asteroid, Hard, Gas, Warp EmiDistib
}

var Opts Options

func init() {
	buf, err := ioutil.ReadFile("emigen_ini.json")
	if err != nil {
		return
		panic(err)
	}
	json.Unmarshal(buf, &Opts)
}
