package main

import (
	"flag"

	"github.com/zycgary/mxshop-go/app/user/service/v1/internal/config"
)

var (
	flagConfig string
)

func init() {
	flag.StringVar(&flagConfig, "config", "app/user/service/v1/configs/config.yaml", "config path, eg: -config config.yaml")
}

func main() {
	flag.Parse()

	// Load config.
	var conf config.Config
	if err := conf.Load(flagConfig); err != nil {
		panic(err)
	}

	a, err := newApp(conf)
	if err != nil {
		panic(err)
	}

	// Start and wait for stop signal.
	if err := a.Run(); err != nil {
		panic(err)
	}
}
