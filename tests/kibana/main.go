package main

import (
	"github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
	"github.com/firstrow/goautosocket"
	"time"
)

func main() {
	log := logrus.New()
	conn, err := gas.Dial("tcp", "magellan2018.aerem.in:5000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{}))

	log.Hooks.Add(hook)

	for {
		rec:=log.WithFields(logrus.Fields{"EventType":"ping event"})
		rec.Infoln("current time ",time.Now())
		time.Sleep(time.Second*15)
	}
}