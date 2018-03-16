//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import (
	"net/http"
	"io"
)

type Server struct{
	mux *http.ServeMux
	httpServ *http.Server
	roomServ RoomGetSetter
}

type RoomGetSetter interface {
	GetRoomCommon (room string) io.Reader
	SetRoomCommon (room string, r io.Reader)
}

const roomAttr = "room"
const roleAttr = "attr"

func roomRole (r *http.Request) (room, role string){
	room = r.Header.Get(roomAttr)
	role = r.Header.Get(roleAttr)
	return
}


func stateHandler(srv *Server) http.Handler {

	f:=func (w http.ResponseWriter, r * http.Request) {
		room, _ := roomRole(r)
		if r.Method == http.MethodPost{
			srv.roomServ.SetRoomCommon(room, r.Body)
		}
	}
	return http.HandlerFunc(f)
}

//NewServer creates a server listening
func NewServer(addr string, roomServ RoomGetSetter) (*Server, error) {
	mux:=http.NewServeMux()
	httpserv := &http.Server{Addr: addr, Handler: mux}
	err:=httpserv.ListenAndServe()
	if err!=nil{
		return nil, err
	}
	srv :=&Server{
		httpServ: httpserv,
		mux:mux,
		roomServ:roomServ,
	}
	mux.Handle("/state/", stateHandler(srv))

	return srv,nil
}
