package main

import (
	"flag"
	"fmt"
	"github.com/washingtonpost/lambda-gateway/server/gateway"
	"log"
	"os"
)

const (
	accessLogTimeFormat = "02/Jan/2006 15:04:05 +0000"
	version             = "1.0.0"
)

var (
	printVersion = flag.Bool("version", false, "Print the version and exit")
	host         = flag.String("host", "unix://tmp/lambda.sock", "Host to listen on. Can be a TCP connection (e.g. tcp://0.0.0.0:8080) or Unix Socket.")
	lambda       = flag.String("lambda", "", "AWS Lambda ARN to invoke (e.g. arn:aws:lambda:us-east-1:111111111111:function:foobar).")
)

func main() {
	flag.Parse()

	if *printVersion {
		fmt.Println("lambda-gateway version", version)
		os.Exit(0)
	}

	start()
}

//sets up REST resources routes and starts the HTTP server
func start() {
	// Initialize loggers
	logger := log.New(os.Stdout, "[server] ", 0)
	engine := gateway.NewEngine(logger, *lambda)
	logger.Fatal(engine.Run(*host))
}
