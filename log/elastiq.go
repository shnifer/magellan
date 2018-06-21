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

/*
used to check elastiq output
type MyTransport struct{
	transport http.RoundTripper
}

func (mt MyTransport) RoundTrip(r *http.Request) (*http.Response, error)  {
	native.Println("======")
	native.Println(r.Header)
	if r.Body!=nil {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			native.Println(err)
		}
		native.Println(string(buf))
	}

	return mt.transport.RoundTrip(r)
}
*/

func Start(timeout, minRetry, maxRetry time.Duration, logTCP string, hostname string) {
	l := logrus.New()
	l.Level = LoggerLevel

	//var myTransport http.RoundTripper = MyTransport{transport: http.DefaultTransport}
	httpClient := &http.Client{
		Timeout: timeout,
	//	Transport: myTransport,
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

		hook, err := elogrus.NewAsyncElasticHook(client, hostname, ELKLevel, indexName)
		if err != nil {
			Log(LVL_ERROR, "elogrus.NewElasticHook", err)
			return
		}
		l.Hooks.Add(hook)
	}

	logger = l
}
