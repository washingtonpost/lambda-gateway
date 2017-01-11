package gateway

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type Engine struct {
	logger     *log.Logger
	lambda_arn string
	config     *aws.Config
}

func NewEngine(logger *log.Logger, lambda_arn string) *Engine {
	creds := credentials.NewEnvCredentials()
	config := &aws.Config{Credentials: creds}
	return &Engine{logger, lambda_arn, config}
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
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	svc := lambda.New(sess)

	params := &lambda.InvokeInput{
		FunctionName: aws.String("arn:aws:lambda:us-east-1:317098396095:function:test-dev-test"),
		Payload:      []byte(`{"message": "hello world"}`),
	}

	resp, err := svc.Invoke(params)

	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%v\n", resp)))
}
