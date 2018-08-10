//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/wrnt"
	"net/http"
	"sync"
	"time"
)

//Options to create and start server
type ServerOpts struct {
	Addr string

	//hooks for server implementation
	RoomServ   RoomCheckGetSetter
	StartState string

	NeededRoles []string

	RoomUpdatePeriod time.Duration
	LastSeenTimeout  time.Duration

	ConsoleHandler func(w http.ResponseWriter, r *http.Request)
}

//Interface all-in-one for server implementation
type RoomCheckGetSetter interface {
	//mb role separation needed, but now Common and State data get full and common
	GetRoomCommon(room string) ([]byte, error)
	SetRoomCommon(room string, data []byte) error
	IsValidState(room string, state string) bool
	RdyStateData(room string, state string)
	GetStateData(room string) []byte

	OnCommand(room string, role string, command string)
}

//network.server data for one room
type servRoomState struct {
	mu sync.Mutex
	//map[role]isOnline
	online   map[string]bool
	lastSeen map[string]time.Time

	state RoomState

	//conns reported new state
	reported map[string]string

	send  *wrnt.SendMany
	recvs map[string]*wrnt.Recv
}

func newServRoomState(commandReceivers []string) *servRoomState {
	defer LogFunc("network.newServRoomState")()

	recvs := make(map[string]*wrnt.Recv)
	for _, name := range commandReceivers {
		recvs[name] = wrnt.NewRecv()
	}

	return &servRoomState{
		online:   make(map[string]bool),
		lastSeen: make(map[string]time.Time),
		reported: make(map[string]string),
		send:     wrnt.NewSendMany(commandReceivers),
		recvs:    recvs,
	}
}

type Server struct {
	mux      *http.ServeMux
	httpServ *http.Server
	opts     ServerOpts

	//Write blocks only to add new room
	mu         sync.RWMutex
	roomsState map[string]*servRoomState

	metric *serverMetric
}

//NewServer creates a server listening
func NewServer(opts ServerOpts) *Server {
	mux := http.NewServeMux()
	httpServ := &http.Server{Addr: opts.Addr, Handler: mux}

	if opts.RoomUpdatePeriod == 0 {
		opts.RoomUpdatePeriod = ServerDefaultRoomUpdatePeriod
	}
	if opts.LastSeenTimeout == 0 {
		opts.LastSeenTimeout = ServerDefaultLastSeenTimeout
	}

	srv := &Server{
		httpServ:   httpServ,
		mux:        mux,
		opts:       opts,
		roomsState: make(map[string]*servRoomState),
	}
	srv.metric = newServerMetric(srv)

	if opts.ConsoleHandler!=nil {
		mux.Handle(consolePattern, consoleHandler(srv))
	}
	mux.Handle(testPattern, testHandler(srv))
	mux.Handle(pingPattern, pingHandler(srv))
	mux.Handle(roomPattern, roomHandler(srv))
	mux.Handle(statePattern, stateHandler(srv))
	go func() {
		err := httpServ.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()

	go serverRoomUpdater(srv)

	return srv
}

//sends command to all clients in room
func (s *Server) AddCommand(roomName string, command string) {
	go func() {
		s.mu.RLock()
		room, ok := s.roomsState[roomName]
		s.mu.RUnlock()
		if !ok {
			return
		}
		room.mu.Lock()
		room.send.AddItems(command)
		room.mu.Unlock()
	}()
}

//used as goroutine
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
		defer LogFunc("network.stateHandler f")()

		srv.metric.add(metricState, metricRPS, 1)

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

		CommonData, err := srv.opts.RoomServ.GetRoomCommon(roomName)
		if err != nil {
			sendErr(w, "can't get fresh Common data to send with state data!")
		}
		SendData := StateDataResp{
			StateData:   srv.opts.RoomServ.GetStateData(roomName),
			StartCommon: CommonData,
		}
		buf := &bytes.Buffer{}
		enc := gob.NewEncoder(buf)
		err = enc.Encode(SendData)
		if err != nil {
			Log(LVL_ERROR, "error: gob.encode.SendData: ", err)
		}

		srv.metric.add(metricState, metricRespBPS, buf.Len())

		buf.WriteTo(w)
	}
	return http.HandlerFunc(f)
}

//room.mu must be already locked
func setNewState(srv *Server, room *servRoomState, roomName, newState string, dropCommands bool) bool {
	if !room.state.IsCoherent {
		Log(LVL_ERROR, "already changing state!", newState)
		return false
	}
	if room.state.Wanted == newState {
		Log(LVL_ERROR, "state is the same")
		return false
	}
	if !srv.opts.RoomServ.IsValidState(roomName, newState) {
		Log(LVL_ERROR, "not valid state", newState)
		return false
	}

	room.state.IsCoherent = false
	room.state.Wanted = newState
	room.state.RdyServData = false

	//flush commands
	if dropCommands {
		room.send.DropNotSent()
	}

	go requestStateData(srv, roomName, newState)
	return true
}

func serverRoomUpdater(serv *Server) {
	defer LogFunc("network.stateHandler")()

	//Only update last seen timeout
	//state update's on changes
	tick := time.Tick(serv.opts.RoomUpdatePeriod)
	for {
		<-tick
		serv.mu.RLock()

		//Conns update
		now := time.Now()
		for _, room := range serv.roomsState {
			room.mu.Lock()
			for role, lastSeen := range room.lastSeen {
				if now.Sub(lastSeen) > serv.opts.LastSeenTimeout {
					room.online[role] = false
				}
			}
			room.mu.Unlock()
		}

		serv.mu.RUnlock()
	}
}

func testHandler(srv *Server) http.Handler {
	_ = srv
	f := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, srv.Metric())
	}

	return http.HandlerFunc(f)
}

func consoleHandler(srv *Server) http.Handler {
	return http.HandlerFunc(srv.opts.ConsoleHandler)
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
	w.Header().Set("error", "1")
	w.Write([]byte(err))
}
