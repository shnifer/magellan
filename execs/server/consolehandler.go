package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"strconv"
)

func (rs *roomServer) consoleHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	url:=r.URL.String()
	log.Println(url)
	parts := strings.Split(url,"/")
	n:=0
	for n<len(parts){
		if parts[n]==""{
			parts = append(parts[:n],parts[n+1:]...)
		} else {
			n++
		}
	}

	if len(parts)<2{
		fmt.Fprintln(w, "use command after /console, i.e. /console/restore/")
		return
	}

	cmd:=parts[1:]
	switch cmd[0]{
	case "restore":
		rs.consoleRestore(cmd, w, r)
	default:
		fmt.Fprintln(w, "unknown command ", cmd[0])
	}
}

func (rs *roomServer) consoleRestore (cmd []string, w http.ResponseWriter, r *http.Request) {
	params:=cmd[1:]

	var id string
	var err error

	if len(params) == 0 {
		rs.consoleRestoreList(w, r)
		return
	} else {
		id=params[0]
	}

	var restoreN int
	if len(params) == 1{
		rs.consoleRestoreShipList(id,w,r)
		return
	} else {
		restoreN,err=strconv.Atoi(params[1])
		if err!=nil{
			fmt.Fprintln(w, "Restore N must be integer, got ",params[1])
			return
		}
	}

	fmt.Fprintln(w,"nice ",id, restoreN)
}

func (rs *roomServer) consoleRestoreList(w http.ResponseWriter, r *http.Request){
	rs.RLock()
	defer rs.RUnlock()

	inFlight:=make(map[string]string)
	for roomName, holder:=range rs.holders{
		state:=holder.getState()
		if (state.StateID!=commons.STATE_cosmo && state.StateID!=commons.STATE_warp) ||
			state.ShipID=="" || state.GalaxyID==""{
			continue
		}
		id:=state.ShipID
		if id == ""{
			continue
		}
		inFlight[roomName] = id
	}

	if len(inFlight)==0{
		fmt.Fprintln(w, "in flight: NONE<br>")
	} else {
		fmt.Fprintln(w, "in flight:<br>")
	}
	for room, flightID:=range inFlight{
		fmt.Fprintf(w, `<a href="./%v"> Room: %v Flight: %v </a><br>`, flightID, room, flightID)
	}
	fmt.Fprintln(w, "or use /restore/#flightID")
}

func (rs *roomServer) consoleRestoreShipList(shipId string, w http.ResponseWriter, r *http.Request) {
	list:=rs.getShipRestoreList(shipId)
	if len(list)==0{
		fmt.Fprintln(w, "Restore points not found for ship ", shipId)
		return
	}
	fmt.Fprintln(w, "for ship ", shipId,":<br>")

	for _,p:=range list{
		fmt.Fprintf(w, `<a href="./%v/%v">%v</a><br>`, shipId, p.restN, p.memo)
	}
}