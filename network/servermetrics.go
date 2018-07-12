package network

import (
	"log"
	"sort"
	"sync"
	"time"
)

type serverMetric struct {
	sync.Mutex
	serv    *Server
	total   ServerMetricMsg
	second  ServerMetricMsg
	current ServerMetricMsg
	msgCh   chan ServerMetricMsg
}

func newServerMetric(s *Server) *serverMetric {
	res := &serverMetric{
		serv:  s,
		msgCh: make(chan ServerMetricMsg, 30),
	}
	go daemonMetric(s, res)
	return res
}

func daemonMetric(s *Server, m *serverMetric) {
	tick := time.Tick(time.Second)
	var msg ServerMetricMsg
	for {
		select {
		case <-tick:
			m.Lock()
			m.second = m.current
			m.current = ServerMetricMsg{}
			m.Unlock()
		case msg = <-m.msgCh:
			m.Lock()
			m.total.Add(msg)
			m.current.Add(msg)
			m.Unlock()
		}
	}
}

func (m *serverMetric) add(reqT, reqF int, val int) {
	req := ServerReqMetric{}
	switch reqF {
	case metricRPS:
		req.RPS += val
	case metricReqBPS:
		req.ReqBPS += val
	case metricRespBPS:
		req.RespBPS += val
	default:
		log.Println(1)
	}
	msg := ServerMetricMsg{}
	switch reqT {
	case metricPing:
		msg.Ping = req
	case metricCommon:
		msg.Common = req
	case metricState:
		msg.State = req
	default:
		log.Println(2)
	}
	m.msgCh <- msg
}

func (m *serverMetric) get() MetricResp {
	m.Lock()
	total := m.total
	second := m.second
	m.Unlock()

	resp := newMetricResp()
	resp.Second = second
	resp.Total = total

	m.serv.mu.RLock()
	resp.RoomCount = len(m.serv.roomsState)

	var totalOnline int
	for roomName, room := range m.serv.roomsState {
		resp.RoomNames = append(resp.RoomNames, roomName)
		room.mu.Lock()
		resp.States[roomName] = room.state.String()
		for _, online := range room.online {
			if online {
				resp.RoomOnline[roomName]++
				totalOnline++
			}
		}
		room.mu.Unlock()
	}
	m.serv.mu.RUnlock()
	resp.OnlineTotal = totalOnline
	sort.Strings(resp.RoomNames)
	return resp
}

func (s *Server) Metric() MetricResp {
	return s.metric.get()
}
