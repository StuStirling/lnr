package main

import (
	"os"

	"github.com/stustirling/lnr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
