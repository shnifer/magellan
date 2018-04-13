//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import (
	"encoding/gob"
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

	NeededRoles []string
}

type RoomCheckGetSetter interface {
	//mb role separation needed, but now Common and State data get full and common
	GetRoomCommon(room string) ([]byte, error)
	SetRoomCommon(room string, data []byte) error
	IsValidState(room string, state string) bool
	RdyStateData(room string, state string)
	GetStateData(room string) []byte

	OnCommand(room string, role string, command string)
}

type servRoomState struct {
	mu sync.Mutex
	//map[role]isOnline
	online   map[string]bool
	lastSeen map[string]time.Time

	state RoomState

	//conns reported new state
	reported map[string]string

	//conns last received Client->Server command
	//map[role]LastCommandReceived (client numeration)
	lastCommandFromClient map[string]int

	//conns last received Server->Client command
	//map[role]LastReceivedCommandN (server numeration)
	lastCommandToClient map[string]int

	//server numeration
	baseCommandN int
	commands     []string
}

func newServRoomState() *servRoomState {
	return &servRoomState{
		online:                make(map[string]bool),
		lastSeen:              make(map[string]time.Time),
		reported:              make(map[string]string),
		lastCommandFromClient: make(map[string]int),
		baseCommandN:          1, //start from #1 server command
		commands:              make([]string, 0),
		lastCommandToClient:   make(map[string]int),
	}
}

type Server struct {
	mux      *http.ServeMux
	httpServ *http.Server
	opts     ServerOpts

	//Write blocks only to add new room
	mu         sync.RWMutex
	roomsState map[string]*servRoomState
}

func (s *Server) AddCommand(roomName string, command string) {
	go func() {
		s.mu.RLock()
		room, ok := s.roomsState[roomName]
		s.mu.RUnlock()
		if !ok {
			return
		}
		room.mu.Lock()
		room.commands = append(room.commands, command)
		room.mu.Unlock()
	}()
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

func stateHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
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

		StateData := srv.opts.RoomServ.GetStateData(roomName)
		CommonData, err := srv.opts.RoomServ.GetRoomCommon(roomName)
		if err != nil {
			sendErr(w, "can't get fresh Common data to send with state data!")
		}
		SendData := StateDataResp{
			StateData:   StateData,
			StartCommon: CommonData,
		}
		enc := gob.NewEncoder(w)
		err = enc.Encode(SendData)
		if err != nil {
			panic(err)
		}
	}
	return http.HandlerFunc(f)
}

func setNewState(srv *Server, room *servRoomState, roomName, newState string) bool {
	if !room.state.IsCoherent {
		log.Println("already changing state!", newState)
		return false
	}
	if room.state.Wanted == newState {
		log.Println("state is the same")
		return false
	}
	if !srv.opts.RoomServ.IsValidState(roomName, newState) {
		log.Println("not valid state", newState)
		return false
	}

	room.state.IsCoherent = false
	room.state.Wanted = newState
	room.state.RdyServData = false

	//flush commands
	room.baseCommandN += len(room.commands)
	room.commands = room.commands[:0]

	go requestStateData(srv, roomName, newState)
	return true
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
func NewServer(opts ServerOpts) *Server {
	mux := http.NewServeMux()
	httpServ := &http.Server{Addr: opts.Addr, Handler: mux}

	srv := &Server{
		httpServ:   httpServ,
		mux:        mux,
		opts:       opts,
		roomsState: make(map[string]*servRoomState),
	}

	mux.Handle(pingPattern, pingHandler(srv))
	mux.Handle(roomPattern, roomHandler(srv))
	mux.Handle(statePattern, stateHandler(srv))
	go func() {
		err := httpServ.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	go serverRoomUpdater(srv)

	return srv
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
