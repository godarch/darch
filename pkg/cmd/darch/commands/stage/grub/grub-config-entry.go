package grub

import (
	"fmt"
	"github.com/godarch/darch/pkg/block"
	"github.com/godarch/darch/pkg/cmd/darch/commands"
	"github.com/godarch/darch/pkg/grub"
	"github.com/urfave/cli"
	"io"
	"os"
)

var grubConfigEntryCommand = cli.Command{
	Name:        "config-entry",
	Description: "outputs grub code to include /etc/darch/grub.cfg on boot",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "cryptodisk,c",
			Usage: "enable cryptodisk feature (for encrypted /boot)",
		},
	},
	Action: func(clicontext *cli.Context) error {
		err := commands.CheckForRoot()
		if err != nil {
			return err
		}

		configPath := "/etc/darch/"

		device, err := block.GetBlockDeviceForPath(configPath)
		if err != nil {
			return err
		}

		relativePathToDevice, err := block.GetPathRelativeToBlockDevice(configPath)
		if err != nil {
			return err
		}

		// Write the required grub code to access the device that our darch grub.cfg exists.
		err = grub.PrepareAccessToDevice(device, os.Stdout, clicontext.Bool("cryptodisk"))
		if err != nil {
			return err
		}

		// Write the code that actually sources our grub.cfg file.
		io.WriteString(os.Stdout, fmt.Sprintf("if [ -f %s/grub.cfg ]; then\n", relativePathToDevice))
		io.WriteString(os.Stdout, fmt.Sprintf("  source %s/grub.cfg\n", relativePathToDevice))
		io.WriteString(os.Stdout, "fi\n")

		return nil
	},
}
