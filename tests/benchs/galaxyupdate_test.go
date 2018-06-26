package main

import (
	"testing"
	."github.com/Shnifer/magellan/commons"
	"io/ioutil"
	"encoding/json"
	"github.com/Shnifer/magellan/v2"
)

var galaxy *Galaxy

func BenchmarkUpdate1(b *testing.B) {
	loadgalaxy()
	var st float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st+=0.1
		Update1(st)
	}
}

func BenchmarkUpdate2(b *testing.B) {
	loadgalaxy()
	var st float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st+=0.1
		Update2(st)
	}
}

func Update1(st float64) {
	for _, obj := range galaxy.Ordered {
		if obj.ParentID == "" {
			continue
		}
		parent := galaxy.Points[obj.ParentID].Pos
		angle := (360 / obj.Period) * st
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}

func Update2(st float64) {
	posMap := make(map[string]v2.V2)

	for _, obj := range galaxy.Ordered {
		if obj.ParentID == "" {
			continue
		}
		parent,ok:=posMap[obj.ParentID]
		if !ok {
			parent = galaxy.Points[obj.ParentID].Pos
			posMap[obj.ParentID] = parent
		}
		angle := (360 / obj.Period) * st
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}

func loadgalaxy(){
	buf,err:=ioutil.ReadFile("testgalaxy.json")
	if err!=nil{
		panic(err)
	}
	err=json.Unmarshal(buf,&galaxy)
	if err!=nil{
		panic(err)
	}

	//First restore ID's
	for id, v := range galaxy.Points {
		if v.ID == "" {
			v.ID = id
			galaxy.Points[id] = v
		}
	}
	//Second - recalc lvls!
	galaxy.RecalcLvls()
}