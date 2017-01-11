FROM golang:1.7-alpine
CMD ["/go/bin/lambda-gateway"]
WORKDIR /go/src/github.com/washingtonpost/lambda-gateway/server
RUN apk update && \
    apk upgrade && \
    apk add --update --no-cache git
RUN go get -u github.com/kardianos/govendor
ADD vendor /go/src/github.com/washingtonpost/lambda-gateway/server/vendor
RUN govendor sync

ADD gateway /go/src/github.com/washingtonpost/lambda-gateway/server/gateway
ADD cmd /go/src/github.com/washingtonpost/lambda-gateway/server/cmd
RUN govendor test +local && \
    go build -o /go/bin/lambda-gateway github.com/washingtonpost/lambda-gateway/server/cmd
