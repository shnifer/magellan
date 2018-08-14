package main

import (
	"github.com/Shnifer/magellan/alice"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"strconv"
	"time"
)

func init() {
	opts := alice.Opts{
		Addr:     DEFVAL.AliceAddr,
		Path:     DEFVAL.AlicePass,
		Password: DEFVAL.AlicePass,
		Login:    DEFVAL.AliceLogin,
	}
	alice.InitAlice(opts)
}

const TryCount = 10
const PauseS = 5

func sendAlice(bioInf, nucleo [7]int) {
	ClientLogGame(Client, "alice", bioInf, nucleo)

	events := make(alice.Events, 0)
	if bioInf != [7]int{} {
		events = append(events, alice.Event{
			EvType: "biological-systems-influence",
			Data:   bioInf,
		})
	}
	if nucleo != [7]int{} {
		events = append(events, alice.Event{
			EvType: "modify-nucleotide-instant",
			Data:   nucleo,
		})
	}
	if len(events) == 0 {
		return
	}

	location := "ship_" + strconv.Itoa(Data.BSP.Dock)
	var err error
	for i := 0; i < TryCount; i++ {
		err = alice.DoReq(location, events)
		if err == nil {
			return
		}
		time.Sleep(PauseS * time.Second)
	}

	LogGame("failedReqs", false,
		"can't send to alice ", err, "req: ", location, events)
}
