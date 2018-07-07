package main

import "github.com/Shnifer/magellan/alice"

func main() {
	opts := alice.Opts{
		Addr:     "https://api.alice.magellan2018.ru",
		Path:     "location_events",
		Password: "admin",
		Login:    "admin",
	}
	alice.InitAlice(opts)

	events := make(alice.Events, 2)
	events[0] = alice.Event{
		EvType: "biological-systems-influence",
		Data:   [7]int{0, 0, 1, 0, 1, 0, 0},
	}
	events[1] = alice.Event{
		EvType: "modify-nucleotide-instant",
		Data:   [7]int{0, -1, 0, 0, 0, 0, 0},
	}

	err := alice.DoReq("ship_3", events)
	if err != nil {
		panic(err)
	}
}
