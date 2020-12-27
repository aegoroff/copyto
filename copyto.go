package main

import (
	"github.com/aegoroff/copyto/cmd"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
