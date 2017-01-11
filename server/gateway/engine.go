package gateway

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Engine struct {
	logger *log.Logger
}

func NewEngine(logger *log.Logger) *Engine {
	return &Engine{logger}
}

func (e *Engine) Run(host string) error {
	if strings.HasPrefix(host, "tcp://") {
		listen := host[6:]
		return e.RunHTTP(listen)
	} else if strings.HasPrefix(host, "unix://") {
		listen := host[6:]
		return e.RunUnix(listen)
	} else {
		return fmt.Errorf("Unable to parse host option")
	}
}

func (e *Engine) RunHTTP(address string) error {
	e.logger.Printf("listening on TCP %s", address)
	return http.ListenAndServe(address, e)
}

func (e *Engine) RunUnix(file string) error {
	e.logger.Printf("listening on UNIX Socket %s", file)
	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		return err
	}
	defer listener.Close()
	return http.Serve(listener, e)
}

// Conforms to the http.Handler interface.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
}
