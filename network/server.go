//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import (
	"io"
	"log"
	"net/http"
	"time"
	"sync"
)

type ServerOpts struct {
	Addr     string

	//hooks for server implementation
	RoomServ RoomCheckGetSetter
}

type RoomCheckGetSetter interface {
	CheckRoomFull(members RoomMembers) bool
	GetRoomCommon(room string) ([]byte, error)
	SetRoomCommon(room string, r io.Reader) error

}


//for implementation of opts.RoomFull()
type RoomMembers map[string]bool

//in room map[role]isConnected
type tRoom struct{
	online map[string]bool
	lastSeen map[string]time.Time
}
func newTRoom() tRoom{
	return tRoom{
		online: make(map[string]bool),
		lastSeen: make(map[string]time.Time),
	}
}

type Server struct {
	mux      *http.ServeMux
	httpServ *http.Server
	opts     ServerOpts

	//map [room]
	//mutex protect the map only
	mu sync.Mutex
	roomsConns map[string]tRoom
}

func (s *Server) checkFullRoom(room tRoom) bool{
	res:=make(RoomMembers)
	for key,val:=range room.online{
		if val{
			res[key] = true
		}
	}

	return s.opts.RoomServ.CheckRoomFull(res)
}

func roomHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		room, _ := roomRole(r)

		//Update room's common state if needed
		if r.Method == POST {
			err := srv.opts.RoomServ.SetRoomCommon(room, r.Body)
			if err != nil {
				log.Println("CANT POST in stateHandler for room", room)
			}
		}

		//Response with new common state
		buf, err := srv.opts.RoomServ.GetRoomCommon(room)
		if err != nil {
			log.Println("CANT GET in stateHandler for room", room)
			return
		} else {
			//do not write response if get failed
			w.Write(buf)
		}
	}
	return http.HandlerFunc(f)
}

func pingHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		srv.mu.Lock()
		defer srv.mu.Unlock()

		roomName, roleName := roomRole(r)

		room,ok:=srv.roomsConns[roomName]
		if !ok{
			room = newTRoom()
			srv.roomsConns[roomName] = room
		}

		room.online[roleName] = true
		room.lastSeen[roleName] = time.Now()

		//Response with room full or not
		isFull:=srv.checkFullRoom(room)

		if isFull{
			w.Write([]byte(MSG_FullRoom))
		} else
		{
			w.Write([]byte(MSG_HalfRoom))
		}
	}
	return http.HandlerFunc(f)
}

func serverRoomUpdater(serv *Server){

	tick:=time.Tick(ServerRoomUpdatePeriod)
	for {
		<-tick
		serv.mu.Lock()
		now:=time.Now()
		for roomName, room := range serv.roomsConns{
			for role, lastSeen := range room.lastSeen{
				if now.Sub(lastSeen)>ServerLastSeenTimeout {
					serv.roomsConns[roomName].online[role] = false
				}
			}
		}
		serv.mu.Unlock()
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
		roomsConns: make(map[string]tRoom),
	}

	mux.Handle(pingPattern, pingHandler(srv))
	mux.Handle(roomPattern, roomHandler(srv))
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