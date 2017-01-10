package main

import (
	"flag"
	"fmt"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"os"
	"strings"
)

const (
	accessLogTimeFormat = "02/Jan/2006 15:04:05 +0000"
	version             = "1.0.0"
)

var (
	printVersion = flag.Bool("version", false, "Print the version and exit")
	//host         = flag.String("host", "tcp://0.0.0.0:8080", "Host to listen on. Can be a TCP connection or Unix Socket.")
	host = flag.String("host", "tcp://0.0.0.0:8080", "Host to listen on. Can be a TCP connection or Unix Socket.")
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
	//accessLogger := log.New(os.Stdout, "[http] ", 0)

	router := gin.Default()
	router.GET("/health", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"status": "ok",
		})
	})
	if strings.HasPrefix(*host, "tcp://") {
		listen := (*host)[6:]
		logger.Printf("listening on TCP %s", listen)
		logger.Fatal(router.Run(listen))
	} else if strings.HasPrefix(*host, "unix://") {
		listen := (*host)[6:]
		logger.Printf("listening on UNIX Socket %s", listen)
		logger.Fatal(router.RunUnix(listen))
	} else {
		logger.Fatal("Unable to parse host option")
	}
}

/*
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

type AccessLogger struct {
	*log.Logger
}

func NewAccessLogger(logger *log.Logger) *AccessLogger {
	return &AccessLogger{logger}
}

func (a *AccessLogger) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	startTime := time.Now()
	next(w, req)

	if req.RequestURI != "/health" { //ignore the health checks from ELB
		res := w.(negroni.ResponseWriter)
		a.Printf("%v %v %v [%v] \"%v %v %v\" %v %v %v",
			clientIP(req),
			"-",
			"-",
			responseTime(),
			req.Method,
			req.RequestURI,
			req.Proto,
			res.Status(),
			res.Size(),
			responseTimeTaken(startTime))
	}
}

func responseTime() string {
	now := time.Now().UTC()
	return now.Format(accessLogTimeFormat)
}

func responseTimeTaken(startTime time.Time) int64 {
	finishTime := time.Now()
	elapsedTime := finishTime.Sub(startTime)
	return elapsedTime.Nanoseconds() / int64(time.Microsecond)
}

func clientIP(req *http.Request) string {
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = req.RemoteAddr
		}
	}

	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}

	return ip
}
*/
