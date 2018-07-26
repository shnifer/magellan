package main

//COPYPASTE IN PILOT

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/storage"
	"strings"
)

func changeCP(objs map[string]*CosmoPoint, id string, point *CosmoPoint) {
	if point == nil {
		Log(LVL_ERROR, "scene change CosmoPoint with nil value")
		return
	}
	if _, ok := objs[id]; ok {
		*objs[id] = *point
	} else {
		objs[id] = point
	}
}

func (s *cosmoScene) addBuilding(b Building) {
	switch b.Type {
	case BUILDING_MINE, BUILDING_FISHHOUSE:
		pd, ok := Data.Galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "cosmoscene addBuilding: can't find added mine or fishhouse on planet", b.PlanetID)
			return
		}
		//update data, saving pointer. We scan CosmoPoint, we don't need new
		changeCP(s.objects, b.PlanetID, NewCosmoPoint(pd, s.cam.Phys()))
	case BUILDING_BEACON, BUILDING_BLACKBOX:
		pd, ok := Data.Galaxy.Points[b.FullKey]
		if !ok {
			Log(LVL_ERROR, "cosmoscene addBuilding: can't find added", b.Type, "with fullkey", b.FullKey)
			return
		}
		//register new object
		changeCP(s.objects, b.FullKey, NewCosmoPoint(pd, s.cam.Phys()))
	default:
		Log(LVL_ERROR, "cosmoscene addBuilding, unknown building type", b.Type)
	}
}

func (s *cosmoScene) delBuilding(b Building) {
	switch b.Type {
	case BUILDING_MINE, BUILDING_FISHHOUSE:
		pd, ok := Data.Galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "cosmoscene delBuilding: can't find added mine on planet", b.PlanetID)
			return
		}
		//update data, saving pointer. We scan CosmoPoint, we don't need new
		changeCP(s.objects, b.PlanetID, NewCosmoPoint(pd, s.cam.Phys()))
	case BUILDING_BEACON, BUILDING_BLACKBOX:
		if _, ok := s.objects[b.FullKey]; !ok {
			Log(LVL_ERROR, "cosmoscene delBuilding: can't del", b.Type, "with fullkey", b.FullKey)
			return
		}
		delete(s.objects, b.FullKey)
	default:
		Log(LVL_ERROR, "cosmoscene delBuilding, unknown building type", b.Type)
	}
}

func (s *cosmoScene) OnCommand(command string) {
	switch {
	case strings.HasPrefix(command, CMD_BUILDINGEVENT):
		buf := []byte(strings.TrimPrefix(command, CMD_BUILDINGEVENT))
		way, b, err := DecodeEvent(buf)
		if err != nil {
			Log(LVL_ERROR, "onCommand can't decode event", string(buf), ":", err)
		}
		if way == storage.Add {
			Data.Galaxy.AddBuilding(b)
			s.addBuilding(b)
		} else if way == storage.Remove {
			Data.Galaxy.DelBuilding(b)
			s.delBuilding(b)
		}
	default:
		Log(LVL_ERROR, "OnCommand strange command:", command)
	}
}
