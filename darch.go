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
	app.Version = Version
	app.Commands = []cli.Command{
		build.Command(),
		extract.Command(),
		stage.Command(),
		pull.Command(),
		push.Command(),
		cli.Command{
			Name:  "version",
			Usage: "Print version information about darch.",
			Action: func(c *cli.Context) error {
				fmt.Printf("version %s\n", Version)
				fmt.Printf("commit %s\n", GitCommit)
				if len(Version) > 0 {
					fmt.Printf("VersionPrerelease %s\n", VersionPrerelease)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
