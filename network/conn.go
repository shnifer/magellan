//package for in game client-server communication
package network

import "net/http"

type Server struct{
	mux *http.ServeMux
	httpserv *http.Server
}

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

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)){
	s.mux.HandleFunc(pattern, handler)
}
