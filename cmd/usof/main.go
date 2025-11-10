package main

import (
	"os"

	"github.com/RofaBR/Go-Usof/cmd/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
