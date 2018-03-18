package main

import (
	"encoding/json"
	"log"
)

type ShipPart struct {
	X, Y    float64
	Caption string
}

type PlanetPart struct {
	W, H int
	X, Y float64
}

type CommonData struct {
	Ship    *ShipPart
	Planets []PlanetPart
}

func main() {

	data := CommonData{
		Ship: &ShipPart{X: 4.4, Y: 5.1, Caption: "firefly"},
		Planets: []PlanetPart{
			PlanetPart{W: 40, H: 40, X: 0.5, Y: 1.5},
			PlanetPart{W: 65, H: 65, X: 12.8, Y: 18.2},
		},
	}
	log.Println(data.Ship, data.Planets)

	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	log.Println(string(b))

	readData := CommonData{}
	err = json.Unmarshal(b, &readData)
	if err != nil {
		panic(err)
	}
	log.Println(readData.Ship, readData.Planets)
}
