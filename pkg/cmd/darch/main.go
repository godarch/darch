package main

import (
	"github.com/godarch/darch/pkg/darchtest"
	"log"
)

// GitCommit The git commit that was compiled. This will be filled in by the compiler.
var GitCommit string

// Version The main version number that is being run at the moment.
var Version = "0.1.0"

func main() {
	err := darchtest.Test()
	if err != nil {
		log.Println(err)
	}
}
