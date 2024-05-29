package main

import (
	"flag"
	"fmt"

	svc "github.com/zycgary/mxshop-go/internal/user/grpc"
	"github.com/zycgary/mxshop-go/internal/user/grpc/config"
)

var (
	env = flag.String("env", "local", "The running environment of the service")
)

func main() {
	flag.Parse()

	// Load config.
	var conf config.Config
	if err := conf.Load("config/user", fmt.Sprintf("grpc.%s", *env), "yaml"); err != nil {
		panic(err)
	}

	a, err := svc.NewApp(conf)
	if err != nil {
		panic(err)
	}

	// Start and wait for stop signal.
	if err := a.Run(); err != nil {
		panic(err)
	}
}
