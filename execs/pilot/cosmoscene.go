package main

import (
	"github.com/hajimehoshi/ebiten"
	."github.com/Shnifer/magellan/commons"
)

type cosmoShip struct{
	pos CShipPos

}

type cosmoScene struct{
	ship *cosmoShip
	objects map[string]*cosmoObject
}

func newCosmoScene() *cosmoScene{
	return &cosmoScene{
		objects: make (map[string]*cosmoObject),
	}
}

func (cosmoScene) Init() {
}

func (cosmoScene) Update(dt float64) {
}

func (cosmoScene) Draw(image *ebiten.Image) {
}

func (cosmoScene) Destroy() {
}


