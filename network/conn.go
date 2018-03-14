//package for in game client-server communication
package network

import (
	"net"
	"bufio"
	"sync"
)

//conn.Recv and Send channels buffer size
const connBufSize = 128

type conn struct{
	c net.Conn
	Recv <-chan string
	Send chan<- string
}

func connListener(conn net.Conn, recv chan string) {
	defer close(recv)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		str := scanner.Text()
		recv <- str
	}
	if err := scanner.Err(); err != nil {
	}
}

func connSender(conn net.Conn, send chan string) {
	writer := bufio.NewWriter(conn)
	for msg := range send {
		if _, err := writer.WriteString(msg + "\n"); err != nil {
			break
		}
		if err := writer.Flush(); err != nil {
			break
		}
	}
}

func newConn(c net.Conn) conn {
	outCh := make(chan string, connBufSize)
	inCh := make(chan string, connBufSize)

	go connListener(c, inCh)
	go connSender(c, outCh)

	res:=conn{
		c: c,
		Recv: inCh,
		Send: outCh,
	}
	return res
}

func (c conn) close() {
	close(c.Send)
	c.c.Close()
}

type Server struct{
	listener net.Listener
	mu sync.Mutex
	conns []conn
}

func NewServer(listener net.Listener) *Server{
	res:=&Server{
		listener:listener,
	}
	go servAccepter(res)
	return res
}

func servAccepter(s *Server) {

	for{
		c,err:=s.listener.Accept()
		if err != nil {
			break
		}

		s.mu.Lock()
		conn:=newConn(c)
		s.conns = append(s.conns, conn)
		s.mu.Unlock()
	}
}

//Close closes all server conns and than close server.listener
func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i:=range s.conns{
		s.conns[i].close()
	}
	s.conns = []conn{}
	s.listener.Close()
}