package main

import (
	"github.com/SailfinIO/agent/pkg/cli"
	"os"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
