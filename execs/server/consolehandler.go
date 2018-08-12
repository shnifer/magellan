package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//localhost:8010/console/...
func (rs *roomServer) consoleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	url := r.URL.String()
	log.Println(url)
	parts := strings.Split(url, "/")
	n := 0
	for n < len(parts) {
		if parts[n] == "" {
			parts = append(parts[:n], parts[n+1:]...)
		} else {
			n++
		}
	}

	if len(parts) < 2 {
		fmt.Fprintln(w, "use command after /console, i.e. /console/restore/")
		return
	}

	cmd := parts[1:]
	switch cmd[0] {
	case "restore":
		rs.consoleRestore(cmd, w, r)
	case "hardkill":
		rs.consoleHardKill(cmd, w, r)
	default:
		fmt.Fprintln(w, "unknown command ", cmd[0])
	}
}

//localhost:8010/console/restore/...
func (rs *roomServer) consoleRestore(cmd []string, w http.ResponseWriter, r *http.Request) {
	params := cmd[1:]

	var id string
	var err error

	//localhost:8010/console/restore/
	if len(params) == 0 {
		rs.consoleRestoreList(w, r)
		return
	} else {
		id = params[0]
	}

	//localhost:8010/console/restore/flightN/
	var restoreN int
	if len(params) == 1 {
		rs.consoleRestoreShipList(id, w, r)
		return
	} else {
		restoreN, err = strconv.Atoi(params[1])
		if err != nil {
			fmt.Fprintln(w, "Restore N must be integer, got ", params[1])
			return
		}
	}

	var roomName string
	//localhost:8010/console/restore/flightN/restoreN/
	if len(params) == 2 {
		rs.consoleRestoreSelectRoom(id, restoreN, w, r)
		return
	} else {
		roomName = params[2]
	}

	fmt.Fprintf(w, "Started process of loading ship %v restore point %v into room %v", id, restoreN, roomName)
	go rs.DoLoadRestore(id, restoreN, roomName)
}

//localhost:8010/console/restore/
func (rs *roomServer) consoleRestoreList(w http.ResponseWriter, r *http.Request) {
	rs.RLock()
	defer rs.RUnlock()

	inFlight := make(map[string]string)
	for roomName, holder := range rs.holders {
		state := holder.getState()
		if (state.StateID != commons.STATE_cosmo && state.StateID != commons.STATE_warp) ||
			state.ShipID == "" || state.GalaxyID == "" {
			continue
		}
		id := state.ShipID
		if id == "" {
			continue
		}
		inFlight[roomName] = id
	}

	if len(inFlight) == 0 {
		fmt.Fprintln(w, "in flight: NONE<br>")
	} else {
		fmt.Fprintln(w, "in flight:<br>")
	}
	for room, flightID := range inFlight {
		fmt.Fprintf(w, `<a href="./%v"> Room: %v Flight: %v </a><br>`, flightID, room, flightID)
	}
	fmt.Fprintln(w, "or use /restore/#flightID")
}

//localhost:8010/console/restore/Ship##/
func (rs *roomServer) consoleRestoreShipList(shipId string, w http.ResponseWriter, r *http.Request) {
	list := rs.getShipRestoreList(shipId)
	if len(list) == 0 {
		fmt.Fprintln(w, "Restore points not found for ship ", shipId)
		return
	}
	fmt.Fprintln(w, "for ship ", shipId, ":<br>")

	for _, p := range list {
		fmt.Fprintf(w, `<a href="./%v/%v">%v</a><br>`, shipId, p.restN, p.memo)
	}
}

//localhost:8010/console/restore/Ship##/Restore##
func (rs *roomServer) consoleRestoreSelectRoom(shipId string, restoreN int, w http.ResponseWriter, r *http.Request) {
	rs.RLock()
	defer rs.RUnlock()

	curRoom := ""
	list := make([]string, 0)

	for roomName, holder := range rs.holders {
		if shipId == holder.getState().ShipID {
			curRoom = roomName
		} else {
			list = append(list, roomName)
		}
	}

	if len(list) == 0 && curRoom == "" {
		fmt.Fprintln(w, "No room found on server")
		return
	}

	if curRoom != "" {
		fmt.Fprintf(w, `Ship %v is already on room %v, <a href="./%v/%v">restore point %v?</a><br>`,
			shipId, curRoom, restoreN, curRoom, restoreN)
	}

	for _, room := range list {
		fmt.Fprintf(w, `<a href="./%v/%v">%v</a><br>`, restoreN, room, room)
	}
}

func (rs *roomServer) consoleHardKill(cmd []string, w http.ResponseWriter, r *http.Request) {
	if len(cmd) < 2 {
		fmt.Fprintln(w, "use hardkill/roomName to drop room without save or anything")
		return
	}

	roomName := cmd[1]
	err := server.KillRoom(roomName)
	if err != nil {
		fmt.Fprintln(w, "error killing room ", roomName, ":", err)
	} else {
		fmt.Fprintln(w, "room ", roomName, "killed and restarted")
	}
}
