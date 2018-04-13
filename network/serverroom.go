package network

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

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

//room.mu must be already locked
func serverReceiveCommands(srv *Server, req CommonReq, room *servRoomState, roomName, roleName string) {

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
			//Do not
			if !room.state.IsCoherent {
				log.Println("STRANGE: COMMAND_CLIENT received while non-coherent. Command: ", command)
				break
			}
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

func serverRecalcCommands(srv *Server, room *servRoomState) {
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
		log.Println("Strange! minimum lastCommandToClient < baseCommandN", delta, "=", minN, "-", room.baseCommandN, "+1 < 0")
	}
	if delta <= 0 {
		return
	}
	if delta > len(room.commands) {
		log.Println("strange! delta>len(commands)", delta, "=", minN, "-", room.baseCommandN, "+1>", len(room.commands))
		delta = len(room.commands)
	}
	room.baseCommandN += delta
	room.commands = room.commands[delta:]
}
