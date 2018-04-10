//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import (
	"encoding/json"
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

type ServRoomState struct {
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

func newServRoomState() *ServRoomState {
	return &ServRoomState{
		online:                make(map[string]bool),
		lastSeen:              make(map[string]time.Time),
		reported:              make(map[string]string),
		lastCommandFromClient: make(map[string]int),
		commands:              make([]string, 0),
		lastCommandToClient:   make(map[string]int),
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
	isFull := true
	for _, roleName := range s.opts.NeededRoles {
		if !room.online[roleName] {
			isFull = false
			break
		}
	}

	room.state.IsFull = isFull
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

		w.Write(srv.opts.RoomServ.GetStateData(roomName))
	}
	return http.HandlerFunc(f)
}

func setNewState(srv *Server, room *ServRoomState, roomName, newState string) bool {
	if !srv.opts.RoomServ.IsValidState(roomName, newState) {
		log.Println("not valid state", newState)
		return false
	}
	if !room.state.IsCoherent {
		log.Println("already changing state!", newState)
		return false
	}
	if room.state.Wanted == newState {
		log.Println("state is the same")
		return false
	}

	room.state.IsCoherent = false
	room.state.Wanted = newState
	room.state.RdyServData = false
	room.baseCommandN += len(room.commands)
	room.commands = room.commands[:0]
	go requestStateData(srv, roomName, newState)
	return true
}

//room.mu must be already locked
func serverReceiveCommands(srv *Server, req CommonReq, room *ServRoomState, roomName, roleName string) {

	alreadyDoneN := room.lastCommandFromClient[roleName]

	for i, command := range req.Commands {
		commandN := req.CommandsBaseN + i
		if commandN <= alreadyDoneN {
			//already done
			continue
		}
		if len(command) < 1 {
			log.Println("empty command!")
			continue
		}
		prefix := command[:1]
		command = command[1:]
		switch prefix {
		case COMMAND_CLIENT:
			room.commands = append(room.commands, command)
			srv.opts.RoomServ.OnCommand(roomName, roleName, command)
		case COMMAND_REQUESTSTATE:
			stateChanged := setNewState(srv, room, roomName, command)
			if stateChanged {
				break
			}
		default:
			log.Println("Strange prefix", prefix)
			continue
		}
	}

	room.lastCommandFromClient[roleName] =
		req.CommandsBaseN + len(req.Commands) - 1
	room.lastCommandToClient[roleName] = req.LastReceivedCommandN
}

func serverRecalcCommands(srv *Server, room *ServRoomState) {
	var minN int
	for _, role := range srv.opts.NeededRoles {
		lastN, ok := room.lastCommandToClient[role]
		if !ok {
			return
		}
		if minN == 0 || lastN < minN {
			minN = lastN
		}
	}
	delta := minN - room.baseCommandN + 1
	if delta < 0 {
		log.Println("Strange! minimum lastCommandToClient < baseCommandN")
	}
	if delta <= 0 {
		return
	}
	if delta > len(room.commands) {
		//TODO: check start state
		log.Println("strange! delta>len(commands)", delta, "=", minN, "-", room.baseCommandN, "+1>", len(room.commands))
		delta = len(room.commands)
	}
	room.baseCommandN += delta
	room.commands = room.commands[delta:]
}

func roomHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		roomName, roleName := roomRole(r)

		reqBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			sendErr(w, "CANT readAll r.body")
			return
		}

		var req CommonReq
		err = json.Unmarshal(reqBuf, &req)
		if err != nil {
			sendErr(w, "CANT unmarshal r.body as CommonReq")
			return
		}

		//Update room's common state if needed
		if req.DataSent {
			err := srv.opts.RoomServ.SetRoomCommon(roomName, []byte(req.Data))
			if err != nil {
				sendErr(w, "CANT POST in stateHandler for room"+roomName+err.Error())
			}
		}

		//Response with new common state
		buf, err := srv.opts.RoomServ.GetRoomCommon(roomName)
		if err != nil {
			sendErr(w, "CANT GET in stateHandler for room"+roomName+err.Error())
			return
		}

		srv.mu.RLock()
		room := srv.roomsState[roomName]
		srv.mu.RUnlock()

		room.mu.Lock()
		defer room.mu.Unlock()

		//receive new commands
		serverReceiveCommands(srv, req, room, roomName, roleName)

		serverRecalcCommands(srv, room)

		resp := CommonResp{
			Data:          string(buf),
			CommandsBaseN: room.baseCommandN,
			Commands:      room.commands,
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			sendErr(w, "can't marshal CommonResp"+err.Error())
		}

		//do not write response if get failed
		w.Write(respBody)
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
			//DOUBLE check cz mutex RUnlock-Lock
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

		pingResp := PingResp{
			Room:                room.state,
			LastCommandReceived: room.lastCommandFromClient[roleName],
		}

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
