//package for in game client-server communication
package network

import (
	"net"
	"bufio"
	"sync"
	"log"
)

//conn.Recv and Send channels buffer size
const connBufSize = 128

type action func()

type conn struct{
	nc   net.Conn
	recv chan string
	send chan string
	onClose action
}

func connListener(c conn) {
	defer close(c.recv)
	scanner := bufio.NewScanner(c.nc)
	for scanner.Scan() {
		str := scanner.Text()
		c.recv <- str
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	c.close()
}

func connSender(c conn) {
	writer := bufio.NewWriter(c.nc)
	for msg := range c.send {
		if _, err := writer.WriteString(msg + "\n"); err != nil {
			break
		}
		if err := writer.Flush(); err != nil {
			break
		}
	}
}

func newConn(c net.Conn, onClose action) conn {
	outCh := make(chan string, connBufSize)
	inCh := make(chan string, connBufSize)

	res:=conn{
		nc:   c,
		recv: inCh,
		send: outCh,
	}
	go connListener(res)
	go connSender(res)

	return res
}

func (c conn) close() {
	close(c.send)
	c.nc.Close()
	if c.onClose!=nil {
		c.onClose()
	}
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
		conn:=newConn(c,nil)
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