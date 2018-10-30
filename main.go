package main

import (
	"mantle/crypt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

func main() {
	if _, err := crypt.Parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
