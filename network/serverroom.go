package network

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"io/ioutil"
	"net/http"
)

func roomHandler(srv *Server) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		srv.metric.add(metricCommon, metricRPS, 1)

		roomName, roleName := roomRole(r)

		reqBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			sendErr(w, "CANT readAll r.body")
			return
		}

		srv.metric.add(metricCommon, metricReqBPS, len(reqBuf))

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
				return
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

		resp := CommonResp{
			Data: string(buf),
		}

		err = room.send.Confirm(roleName, req.ClientConfirmN)
		if err != nil {
			sendErr(w, err.Error())
		}

		message, err := room.send.Pack(roleName)
		if err == nil {
			resp.Message = message
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			sendErr(w, "can't marshal CommonResp"+err.Error())
		}

		srv.metric.add(metricCommon, metricRespBPS, len(respBody))

		//do not write response if get failed
		w.Write(respBody)
	}
	return http.HandlerFunc(f)
}

//room.mu must be already locked
func serverReceiveCommands(srv *Server, req CommonReq, room *servRoomState, roomName, roleName string) {

	if !room.state.IsCoherent {
		if len(req.Message.Items) > 0 {
			Log(LVL_WARN, "STRANGE: COMMANDs recieved while non-coherent.")
		}

		return
	}

	//we do NOT recv commands if room is not coherent
	//so we do not count them as received
	var commands []string
	if rcv, ok := room.recvs[roleName]; ok {
		commands = rcv.Unpack(req.Message)
	} else {
		Log(LVL_ERROR, "serverReceiveCommands unknown roleName ", roleName)
		return
	}

	for _, command := range commands {

		if len(command) < 1 {
			Log(LVL_WARN, "empty command!")
			continue
		}

		prefix := command[:1]
		command := command[1:]
		switch prefix {
		case COMMAND_CLIENTREQUEST:
			srv.opts.RoomServ.OnCommand(roomName, roleName, command)
		case COMMAND_ROOMBROADCAST:
			room.send.AddItems(command)
		case COMMAND_REQUESTSTATE:
			dropCommands := command[:1] == "+"
			command := command[1:]
			setNewState(srv, room, roomName, command, dropCommands)
			//non valid states reported by implementation
			/*stateChanged :=
			if !stateChanged {

			}*/
		default:
			Log(LVL_ERROR, "Strange prefix", prefix)
			continue
		}
	}
}
