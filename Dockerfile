FROM golang:1.7-alpine
CMD ["/go/bin/lambda-gateway"]
WORKDIR /go/src/github.com/washingtonpost/lambda-gateway
RUN apk update && \
    apk upgrade && \
    apk add --update --no-cache git
RUN go get -u github.com/kardianos/govendor
ADD vendor vendor
RUN govendor sync

ADD cmd cmd
ADD gateway gateway
RUN govendor test +local && \
    go build -o /go/bin/lambda-gateway github.com/washingtonpost/lambda-gateway/cmd
