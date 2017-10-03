package main

import (
	"fmt"
	"log"
	"os"

	"./commands/build"
	"./commands/extract"
	"./commands/pull"
	"./commands/push"
	"./commands/stage"
	"github.com/urfave/cli"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func main() {

	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	app := cli.NewApp()
	app.Name = "darch"
	app.Usage = "A tool used to build, boot and share stateless Arch images."

	app.Commands = []cli.Command{
		build.Command(),
		extract.Command(),
		stage.Command(),
		pull.Command(),
		push.Command(),
	}

	app.Run(os.Args)
}
