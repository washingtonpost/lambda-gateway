package main

import (
	"flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	accessLogTimeFormat = "02/Jan/2006 15:04:05 +0000"
	version             = "1.0.0"
)

var (
	printVersion = flag.Bool("version", false, "Print the version and exit")
	listen       = flag.String("listen", "0.0.0.0:8080", "IP and port to listen for HTTP requests (i.e. 0.0.0.0:8080)")
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
	accessLogger := log.New(os.Stdout, "[http] ", 0)

	mux := mux.NewRouter()
	mux.HandleFunc("/health", healthCheckHandler).Methods("GET")

	//start listening for HTTP requests
	n := negroni.New(negroni.NewRecovery(), NewAccessLogger(accessLogger))
	n.UseHandler(mux)
	logger.Printf("listening on %s", *listen)
	logger.Fatal(http.ListenAndServe(*listen, n))
}

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
