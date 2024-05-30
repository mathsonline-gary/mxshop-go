package main

import (
	"flag"

	svc "github.com/zycgary/mxshop-go/internal/user/grpc"
	"github.com/zycgary/mxshop-go/internal/user/grpc/config"
)

var (
	flagConfig string
)

func init() {
	flag.StringVar(&flagConfig, "config", "config/user/grpc/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// Load config.
	var conf config.Config
	if err := conf.Load(flagConfig); err != nil {
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
