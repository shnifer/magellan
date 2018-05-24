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

		room.send.Confirm(roleName, req.ClientConfirmN)

		message, err := room.send.Pack(roleName)
		if err == nil {
			resp.Message = message
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

	commands := room.recvs[roleName].Unpack(req.Message)

	for _, command := range commands {

		if len(command) < 1 {
			log.Println("empty command!")
			continue
		}

		prefix := command[:1]
		command := command[1:]
		switch prefix {
		case COMMAND_CLIENT:
			//ignore commands sent on not coherent state
			if !room.state.IsCoherent {
				log.Println("STRANGE: COMMAND_CLIENT received while non-coherent. Command: ", command)
				continue
			}
			room.send.AddItems(command)
			srv.opts.RoomServ.OnCommand(roomName, roleName, command)
		case COMMAND_REQUESTSTATE:
			//ignore commands sent on not coherent state
			if !room.state.IsCoherent {
				log.Println("Request state command in non coherent room. Is this good?", command)
			}

			stateChanged := setNewState(srv, room, roomName, command)
			if stateChanged {
				break
			}
		default:
			log.Println("Strange prefix", prefix)
			continue
		}
	}
}
