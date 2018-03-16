//package for in game client-server communication using http-requests
//it has a lot of hard assumptions like room-role structure and is not supposed to be universal
//req.header.values are used to pass room-role with each request, both GET and POST

//clients must make send request eventually or be considered disconnected
//package implements disconnect notification for all other clients in the same room
package network

import "net/http"

type Server struct{
	mux *http.ServeMux
	httpserv *http.Server
}

//NewServer creates a server listening
func NewServer(addr string) (*Server, error) {
	mux:=http.NewServeMux()
	httpserv := &http.Server{Addr: addr, Handler: mux}
	err:=httpserv.ListenAndServe()
	if err!=nil{
		return nil, err
	}
	return &Server{
		httpserv: httpserv,
		mux:mux,
	},nil
}

func (s *Server) handleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)){
	s.mux.HandleFunc(pattern, handler)
}
