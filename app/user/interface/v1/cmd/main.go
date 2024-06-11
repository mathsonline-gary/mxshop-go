package main

import (
	"flag"

	"github.com/zycgary/mxshop-go/app/user/interface/v1/internal/config"
)

var (
	flagConfig string
)

func init() {
	flag.StringVar(&flagConfig, "config", "app/user/interface/v1/configs/config.yaml", "config path, eg: -config config.yaml")
}

func main() {
	flag.Parse()

	// Load config.
	var conf config.Config
	if err := conf.Load(flagConfig); err != nil {
		panic(err)
	}

	a, cleanup, err := newApp(&conf)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Start and wait for stop signal.
	if err := a.Run(); err != nil {
		panic(err)
	}
}
