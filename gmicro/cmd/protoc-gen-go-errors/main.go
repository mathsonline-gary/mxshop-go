package main

import (
	"flag"

	gen2 "github.com/zycgary/mxshop-go/gmicro/cmd/protoc-gen-go-errors/gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	flag.Parse()
	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			gen2.GenerateFile(gen, f)
		}
		return nil
	})
}
