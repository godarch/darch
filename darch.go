package main

import (
	"fmt"
	"log"
	"os"

	"./commands/build"
	"./commands/builddep"
	"./commands/hooks"
	"./commands/inspect"
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
	app.Version = Version
	app.HideVersion = true
	app.Commands = []cli.Command{
		build.Command(),
		inspect.Command(),
		stage.Command(),
		builddep.Command(),
		hooks.Command(),
		cli.Command{
			Name:  "version",
			Usage: "Print version information about darch.",
			Action: func(c *cli.Context) error {
				fmt.Printf("version %s\n", Version)
				fmt.Printf("commit %s\n", GitCommit)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
