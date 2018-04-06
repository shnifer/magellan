//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServerOpts struct {
	Addr string

	//hooks for server implementation
	RoomServ   RoomCheckGetSetter
	StartState string
}

type RoomCheckGetSetter interface {
	CheckRoomFull(members RoomMembers) bool

	//mb role separation needed, but now Common and State data get full and common
	GetRoomCommon(room string) ([]byte, error)
	SetRoomCommon(room string, r io.Reader) error
	IsValidState(room string, state string) bool
	RdyStateData(room string, state string)
	GetStateData(room string) []byte
}

//for implementation of opts.RoomFull()
type RoomMembers map[string]bool

type ServRoomState struct {
	mu sync.Mutex
	//map[role]isOnline
	online   map[string]bool
	lastSeen map[string]time.Time

	state RoomState

	//conns reported new state
	reported map[string]string
}

func newServRoomState() *ServRoomState {
	return &ServRoomState{
		online:   make(map[string]bool),
		lastSeen: make(map[string]time.Time),
		reported: make(map[string]string),
	}
}

//update coherency
//do no set Mutex, must be called within critical section
func (r *ServRoomState) updateState() {
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

type Server struct {
	mux      *http.ServeMux
	httpServ *http.Server
	opts     ServerOpts

	//Write blocks only to add new room
	mu         sync.RWMutex
	roomsState map[string]*ServRoomState
}

func (s *Server) checkFullRoom(room *ServRoomState) {
	res := make(RoomMembers)
	for key, val := range room.online {
		if val {
			res[key] = true
		}
	}

	room.state.IsFull = s.opts.RoomServ.CheckRoomFull(res)
}

func requestStateData(srv *Server, roomName string, newState string) {
	//may be implemented by a long time operation, timeouts provided by implementation
	//while GetStateData do not return room state can not be coherented
	srv.opts.RoomServ.RdyStateData(roomName, newState)

	srv.mu.RLock()
	room := srv.roomsState[roomName]
	room.mu.Lock()

	room.state.RdyServData = true

	room.mu.Unlock()
	srv.mu.RUnlock()
}

func postStateHandler(srv *Server, w http.ResponseWriter, r *http.Request) {
	roomName, _ := roomRole(r)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErr(w, "CANT read state request body!")
		return
	}
	newState := string(b)

	if !srv.opts.RoomServ.IsValidState(roomName, newState) {
		sendErr(w, "state is not valid "+newState)
		return
	}

	srv.mu.RLock()
	defer srv.mu.RUnlock()

	room := srv.roomsState[roomName]
	room.mu.Lock()
	defer room.mu.Unlock()

	if !room.state.IsCoherent {
		sendErr(w, "already changing state!")
		return
	}

	if room.state.Wanted == newState {
		sendErr(w, "state is the same")
		return
	}

	room.state.IsCoherent = false
	room.state.Wanted = newState
	room.state.RdyServData = false
	go requestStateData(srv, roomName, newState)
}

func getStateHandler(srv *Server, w http.ResponseWriter, r *http.Request) {
	roomName, _ := roomRole(r)
	srv.mu.RLock()
	defer srv.mu.RUnlock()

	room := srv.roomsState[roomName]
	room.mu.Lock()
	if !room.state.RdyServData {
		sendErr(w, "Serv state Data is not Ready")
		room.mu.Unlock()
		return
	}
	room.mu.Unlock()

	w.Write(srv.opts.RoomServ.GetStateData(roomName))
}

func stateHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case POST:
			//set NEW state for room
			postStateHandler(srv, w, r)
		case GET:
			//get state Data
			getStateHandler(srv, w, r)
		}
	}
	return http.HandlerFunc(f)
}

func roomHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		room, _ := roomRole(r)

		//Update room's common state if needed
		if r.Method == POST {
			err := srv.opts.RoomServ.SetRoomCommon(room, r.Body)
			if err != nil {
				sendErr(w, "CANT POST in stateHandler for room"+room+err.Error())
			}
		}

		//Response with new common state
		buf, err := srv.opts.RoomServ.GetRoomCommon(room)
		if err != nil {
			sendErr(w, "CANT GET in stateHandler for room"+room+err.Error())
			return
		}

		//do not write response if get failed
		w.Write(buf)
	}
	return http.HandlerFunc(f)
}

func pingHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		srv.mu.RLock()
		defer srv.mu.RUnlock()

		clientState := r.Header.Get(stateAttr)
		roomName, roleName := roomRole(r)

		room, ok := srv.roomsState[roomName]
		if !ok {
			srv.mu.RUnlock()
			srv.mu.Lock()
			room, ok = srv.roomsState[roomName]
			if !ok {
				room = newServRoomState()
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

		pingResp := room.state
		b, err := json.Marshal(pingResp)
		if err != nil {
			panic(err)
		}

		w.Write(b)
	}
	return http.HandlerFunc(f)
}

func serverRoomUpdater(serv *Server) {

	//Only update last seen timeout
	//state update's on changes
	tick := time.Tick(ServerRoomUpdatePeriod)
	for {
		<-tick
		serv.mu.RLock()

		//Conns update
		now := time.Now()
		for _, room := range serv.roomsState {
			room.mu.Lock()
			for role, lastSeen := range room.lastSeen {
				if now.Sub(lastSeen) > ServerLastSeenTimeout {
					room.online[role] = false
				}
			}
			room.mu.Unlock()
		}

		serv.mu.RUnlock()
	}
}

//NewServer creates a server listening
func NewServer(opts ServerOpts) (*Server, error) {
	mux := http.NewServeMux()
	httpServ := &http.Server{Addr: opts.Addr, Handler: mux}

	srv := &Server{
		httpServ:   httpServ,
		mux:        mux,
		opts:       opts,
		roomsState: make(map[string]*ServRoomState),
	}

	mux.Handle(pingPattern, pingHandler(srv))
	mux.Handle(roomPattern, roomHandler(srv))
	mux.Handle(statePattern, stateHandler(srv))
	go httpServ.ListenAndServe()

	go serverRoomUpdater(srv)

	return srv, nil
}

func (s *Server) Close() error {
	return s.httpServ.Close()
}

func roomRole(r *http.Request) (room, role string) {
	room = r.Header.Get(roomAttr)
	role = r.Header.Get(roleAttr)
	return
}

func sendErr(w http.ResponseWriter, err string) {
	log.Println("Error : ", err)
	w.Header().Set("error", "1")
	w.Write([]byte(err))
}
