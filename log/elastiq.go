package log

import (
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"net/http"
	"time"
)

var logger *logrus.Logger

const (
	DefaultLevel = logrus.InfoLevel
	indexName    = "logstash"
)

func Start(timeout, minRetry, maxRetry time.Duration, logTCP string) {
	l := logrus.New()
	httpClient := &http.Client{
		Timeout: timeout,
	}

	backoff := elastic.NewExponentialBackoff(minRetry, maxRetry)
	retrier := elastic.NewBackoffRetrier(backoff)

	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL(logTCP),
		elastic.SetHttpClient(httpClient),
		elastic.SetRetrier(retrier),
	)

	if err != nil {
		Log(LVL_ERROR, "elastic.NewClient", err)
		return
	}

	hook, err := elogrus.NewElasticHook(client, "localhost", DefaultLevel, indexName)
	if err != nil {
		Log(LVL_ERROR, "elogrus.NewElasticHook", err)
		return
	}
	l.Hooks.Add(hook)

	logger = l
}
