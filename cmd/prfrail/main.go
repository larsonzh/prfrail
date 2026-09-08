package main

import (
	"context"
	"os"

	"github.com/larsonzh/prfrail/internal/console"
)

var version = "0.1.0-dev"

func main() {
	cli := console.NewCLI(version)
	exitCode := cli.Execute(context.Background(), os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
