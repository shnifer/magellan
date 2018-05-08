package network

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func pingHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(2) == 0 {
			//return
		}
		srv.mu.RLock()
		defer srv.mu.RUnlock()

		clientState := r.Header.Get(stateAttr)
		roomName, roleName := roomRole(r)

		room, ok := srv.roomsState[roomName]
		if !ok {
			srv.mu.RUnlock()
			srv.mu.Lock()
			//DOUBLE check cz mutex RUnlock-Lock
			room, ok = srv.roomsState[roomName]
			if !ok {
				room = newServRoomState(srv.opts.NeededRoles)
				room.state.Wanted = srv.opts.StartState
				go requestStateData(srv, roomName, room.state.Wanted)
				srv.roomsState[roomName] = room
			}
			srv.mu.Unlock()
			srv.mu.RLock()
		}

		room.mu.Lock()
		defer room.mu.Unlock()

		room.online[roleName] = true
		room.lastSeen[roleName] = time.Now()
		srv.checkFullRoom(room)

		room.reported[roleName] = clientState
		room.updateState()

		var lastRecv int
		if recv, ok := room.recvs[roleName]; ok {
			lastRecv = recv.LastRecv()
		} else {
			log.Println("NO RECIEVER ROOM", roomName, "ROLE", roleName)
		}

		pingResp := PingResp{
			Room:           room.state,
			ServerConfirmN: lastRecv,
		}

		b, err := json.Marshal(pingResp)
		if err != nil {
			panic(err)
		}

		w.Write(b)
	}
	return http.HandlerFunc(f)
}

func (s *Server) checkFullRoom(room *servRoomState) {
	isFull := true
	for _, roleName := range s.opts.NeededRoles {
		if !room.online[roleName] {
			isFull = false
			break
		}
	}

	room.state.IsFull = isFull
}

//update coherency
//do no set Mutex, must be called within critical section
func (r *servRoomState) updateState() {
	if !r.state.IsCoherent {
		coherent := true
		for _, state := range r.reported {
			if state != r.state.Wanted {
				coherent = false
				break
			}
		}
		r.state.IsCoherent = coherent
	}
}
