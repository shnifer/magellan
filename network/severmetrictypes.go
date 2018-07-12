package network

import (
	"fmt"
	"strconv"
)

type ServerReqMetric struct {
	//requests per second
	RPS int
	//request.body bytes per second
	ReqBPS int
	//responce.body bytes per second
	RespBPS int
}

const (
	metricPing = iota
	metricState
	metricCommon
)

const (
	metricRPS = iota
	metricReqBPS
	metricRespBPS
)

type ServerMetricMsg struct {
	Ping, State, Common ServerReqMetric
}

func (a *ServerReqMetric) Add(b ServerReqMetric) {
	a.RPS += b.RPS
	a.ReqBPS += b.ReqBPS
	a.RespBPS += b.RespBPS
}

func (a *ServerMetricMsg) Add(b ServerMetricMsg) {
	a.Common.Add(b.Common)
	a.Ping.Add(b.Ping)
	a.State.Add(b.State)
}

type MetricResp struct {
	RoomCount   int
	RoomNames   []string
	OnlineTotal int
	RoomOnline  map[string]int
	States      map[string]string

	Second ServerMetricMsg
	Total  ServerMetricMsg
}

func newMetricResp() MetricResp {
	return MetricResp{
		RoomNames:  make([]string, 0),
		RoomOnline: make(map[string]int),
		States:     make(map[string]string),
	}
}

func (m MetricResp) String() string {
	roomNames := "{"
	roomOnline := "{"
	roomStates := "{"
	for i, roomName := range m.RoomNames {
		if i > 0 {
			roomNames += ", "
			roomOnline += ", "
		}
		roomNames += `"` + roomName + `"`
		roomOnline += `"` + roomName + `": ` + strconv.Itoa(m.RoomOnline[roomName])
		roomStates += `"` + roomName + `": ` + m.States[roomName]
	}
	roomNames += "}"
	roomOnline += "}"
	roomStates += "}"
	frmt := "{RoomCount: %v, OnlineTotal: %v,\nRoomNames: %v, " +
		"RoomOnline: %v,\nSecond: %v,\nTotal: %v\nStates:%v}"
	return fmt.Sprintf(frmt, m.RoomCount, m.OnlineTotal,
		roomNames, roomOnline, m.Second, m.Total, roomStates)
}

func (m ServerMetricMsg) String() string {
	frmt := `{Ping:%v, State:%v, Common:%v}`
	return fmt.Sprintf(frmt, m.Ping, m.State, m.Common)
}

func (m ServerReqMetric) String() string {
	frmt := `{RPS:%v, ReqBPS:%v, RespBPS:%v}`
	return fmt.Sprintf(frmt, m.RPS, m.ReqBPS, m.RespBPS)
}
