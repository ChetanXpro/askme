package main

import (
	"askme/cmd"
	"askme/initialization"

	"github.com/spf13/pflag"
)

var config *initialization.Config

func main() {

	var (
		cfg = pflag.StringP("config", "c", "./config.yaml", "config file path")
	)
	pflag.Parse()
	config = initialization.LoadConfig(*cfg)

	cmd.Execute()
}
