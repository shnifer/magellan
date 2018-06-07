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
	LoggerLevel = logrus.InfoLevel
	ELKLevel    = logrus.InfoLevel
	indexName   = "logstash"
)

func Start(timeout, minRetry, maxRetry time.Duration, logTCP string) {
	l := logrus.New()
	l.Level = LoggerLevel

	httpClient := &http.Client{
		Timeout: timeout,
	}

	if logTCP != "" {
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

		hook, err := elogrus.NewAsyncElasticHook(client, "localhost", ELKLevel, indexName)
		if err != nil {
			Log(LVL_ERROR, "elogrus.NewElasticHook", err)
			return
		}
		l.Hooks.Add(hook)
	}

	logger = l
}
