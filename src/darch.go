package main

import (
	"os"

	"./commands/build"
	"./commands/extract"
	"./commands/pull"
	"./commands/push"
	"./commands/stage"
	"github.com/urfave/cli"
)

func main() {
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
