FROM golang:1.10

WORKDIR /go/src/github.com/Shnifer/magellan

# Build dependencies
RUN go get bytes
RUN go get encoding/json
RUN go get errors
RUN go get github.com/peterbourgon/diskv
RUN go get io/ioutil
RUN go get os
RUN go get os/signal
RUN go get strconv
RUN go get strings
RUN go get sync
RUN go get time
RUN go get github.com/olivere/elastic
RUN go get gopkg.in/sohlich/elogrus.v3
RUN go get github.com/sirupsen/logrus
RUN go get golang.org/x/image/colornames

# Tool do do something before build
RUN go get -u github.com/gobuffalo/packr/...

COPY . .
RUN /go/bin/packr clean
RUN /go/bin/packr

RUN go install ./execs/server/

CMD ["/go/bin/server"]