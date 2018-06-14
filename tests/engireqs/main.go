package main

import (
	"net/http"
	"strconv"
	"log"
	"io/ioutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/Shnifer/magellan/graph"
	"image/color"
	"github.com/Shnifer/magellan/v2"
	"sync"
	"time"
	"encoding/json"
	"strings"
	"fmt"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/Shnifer/magellan/draw"
)

const addr = "http://tagunil.ru:8000/system/"
var slotSprite *graph.Sprite
var dotSprite *graph.Sprite
var smokeSprite *graph.Sprite
type system struct{
	dat uint16
}
type ReqResp struct{
	Programmed uint16 `json:"programmed"`
	Corrected uint16 `json:"corrected"`
	Id uint `json:"timestamp"`
}

var mu sync.Mutex
var corrected [8] system
var programmed[8] system

func run(image *ebiten.Image) error {
	mu.Lock()
	defer mu.Unlock()

	procClick()

	if ebiten.IsRunningSlowly(){return nil}
	image.Fill(color.White)
	drawSlots(image)
	return nil
}

func main(){
	draw.InitTexAtlas()

	slotSprite=draw.NewAtlasSprite("BUILDING_BEACON",graph.NoCam)
	slotSprite.SetSize(50,50)

	dotSprite = draw.NewAtlasSprite("trail",graph.NoCam)
	dotSprite.SetSize(60,60)

	smokeSprite = draw.NewAtlasSprite("smoke",graph.NoCam)
	smokeSprite.SetSize(60,60)

	go func(){
		for {
			time.Sleep(time.Second)
			reqState()
		}
	}()

	for i:=0; i<8;i++{
		sendProgrammed(i)
	}

	ebiten.Run(run, 1325,725,1, "Engi")
}

func reqState() {
	var r ReqResp
	for i:=0; i<8; i++ {
		addr:=addr+strconv.Itoa(i)
		resp,err:=http.Get(addr)
		if err!=nil{
			log.Println("ERROR get request ",addr,":",err)
		}
		defer resp.Body.Close()
		buf,err:=ioutil.ReadAll(resp.Body)
		if err!=nil{
			log.Println("ERROR read body ",addr,":",err)
		}
		json.Unmarshal(buf, &r)
		mu.Lock()
		corrected[i].Set(r.Corrected)
		mu.Unlock()
	}
}

func sendProgrammed(sn int){
	addr:=addr+strconv.Itoa(sn)
	req:=fmt.Sprintf("{\"programmed\":%v}",programmed[sn].dat)
	log.Println("POST",addr,req)
	r:=strings.NewReader(req)
	resp,err:=http.Post(addr,"application/json",r)
	if err!=nil{
		panic(err)
	}
	defer resp.Body.Close()
}

func drawSlots(image *ebiten.Image) {
	for n:=0; n<8;n++{
		for m:=0; m<16; m++{
			p := pos(n,m)
			if programmed[n].GetByte(m) {
				smokeSprite.SetPos(p)
				smokeSprite.Draw(image)
			}
			slotSprite.SetPos(p)
			slotSprite.Draw(image)
			if corrected[n].GetByte(m) {
				dotSprite.SetPos(p)
				dotSprite.Draw(image)
			}
		}
	}
}

func procClick(){
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x,y:=ebiten.CursorPosition()
		v:=v2.V2{X:float64(x),Y:float64(y)}
		for sn:=0;sn<8;sn++{
			for bn:=0;bn<16;bn++{
				d:=pos(sn,bn).Sub(v).Len()
				if d<50{
					log.Println("click",sn,bn)
					programmed[sn].XorByte(bn)
					go sendProgrammed(sn)
				}
			}
		}
	}
}

func pos(s, b int ) v2.V2{
	return v2.V2{100,100}.AddMul(v2.V2{75,0},float64(b)).AddMul(v2.V2{0,75},float64(s))
}

func (s *system) GetByte(b int) bool{
	n:=s.dat>>uint(b)
	return n&1 > 0
}

func (s *system)Set(x uint16) {
	s.dat = x
}

func (s *system)XorByte(b int) {
	s.dat=s.dat^(1<<uint(b))
}