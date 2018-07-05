package main

import (
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
)

func main() {
	log := logrus.New()
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("https://elasticsearch.magellan2018.aerem.in"))
	if err != nil {
		log.Panic(err)
	}
	hook, err := elogrus.NewElasticHook(client, "localhost", logrus.DebugLevel, "logstash")
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)

	log.WithFields(logrus.Fields{
		"name": "joe",
		"age":  42,
	}).Error("Hello world!")
}
