package stage

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"../../hooks"
	"../../images"
	"../../stage"
	"../../utils"
	"github.com/gobwas/glob"
	"github.com/kennygrant/sanitize"
	"github.com/ryanuber/columnize"
	"github.com/urfave/cli"
)

func uploadCommand() cli.Command {
	return cli.Command{
		Name:      "upload",
		Usage:     "Upload an image to the stage to be booted.",
		ArgsUsage: "IMAGE_NAME",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "tag, t",
				Usage: "The tag to stage.",
				Value: "local",
			},
		},
		Action: func(c *cli.Context) error {
			err := upload(c.Args().First(), c.String("tag"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func listCommand() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "List the images currently staged.",
		Action: func(c *cli.Context) error {
			err := list()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func syncBootLoaderCommand() cli.Command {
	return cli.Command{
		Name:  "sync-boot-loader",
		Usage: "Update boot loader to reflect newly created/removed images.",
		Action: func(c *cli.Context) error {
			err := syncBootLoader()
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

func runHooksCommand() cli.Command {
	return cli.Command{
		Name:  "run-hooks",
		Usage: "Run all the hooks for every images.",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "only-images, i",
				Usage: "A globbing pattern matching \"image:tag\" to update hooks for.",
				Value: "",
			},
		},
		Action: func(c *cli.Context) error {
			err := runHooks(c.String("only-images"))
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			return nil
		},
	}
}

// Command Returns the command to be passed to a cli context.
func Command() cli.Command {
	return cli.Command{
		Name:  "stage",
		Usage: "Commands that help manage the stage.",
		Subcommands: []cli.Command{
			uploadCommand(),
			listCommand(),
			syncBootLoaderCommand(),
			runHooksCommand(),
		},
	}
}

func upload(name string, tag string) error {

	if len(name) == 0 {
		return fmt.Errorf("Name is required")
	}

	if len(tag) == 0 {
		return fmt.Errorf("Tag is required")
	}

	destinationDirectory := "/var/darch/staged"
	destinationDirectory = path.Join(destinationDirectory, sanitize.Path(name+"/"+tag))

	log.Println("Name: " + name)
	log.Println("Tag: " + tag)
	log.Println("Destination: " + destinationDirectory)

	err := images.ExtractImage(name, tag, destinationDirectory)

	if err != nil {
		return err
	}

	return runHooks(name + ":" + tag)
}

func list() error {
	stagedItems, err := stage.GetAllStaged("/var/darch/staged")

	if err != nil {
		return err
	}

	result := []string{
		"Name | Tag | Path | Kernel | InitramFS | RootFS",
	}

	for _, stagedItem := range stagedItems {
		for _, stagedItemTag := range stagedItem.Tags {
			result = append(result, stagedItem.Name+" | "+stagedItemTag.Name+" | "+stagedItemTag.Path+" | "+stagedItemTag.BootKernel+" | "+stagedItemTag.BootInitRAMFS+" | "+stagedItemTag.BootRootFS)
		}
	}

	fmt.Println(columnize.SimpleFormat(result))

	return nil
}

func syncBootLoader() error {
	// It may seem weird that we are just wrapping "grub-mkconfig".
	// 1) I darch-related operations to be done through a single cli interface.
	// 2) The process of updating bootloaders may change from OS versions.
	//    Recommending people go through this method gives me a point at which
	//    I can add if-logic for different os/image types.
	if !utils.FileExists("/etc/grub.d/60_darch") {
		return fmt.Errorf("Grub generator doesn't exist at %s", "/etc/grub.d/60_darch")
	}

	log.Println("Generating grub boot entries...")

	cmd := exec.Command("grub-mkconfig", "--output", "/boot/grub/grub.cfg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runHooks(onlyImages string) error {
	stagedItems, err := stage.GetAllStaged("/var/darch/staged")

	if err != nil {
		return err
	}

	allHooks, err := hooks.GetHooks()

	if err != nil {
		return err
	}

	for _, stagedItem := range stagedItems {
		for _, stagedItemTag := range stagedItem.Tags {
			allowed := true
			if len(onlyImages) > 0 {
				// Let's make sure the glob matches this image
				g := glob.MustCompile(onlyImages)
				if !g.Match(stagedItemTag.FullName) {
					allowed = false
				}
			}
			if allowed {
				// First, let's delet all hold holds
				err = os.RemoveAll(path.Join(stagedItemTag.Path, "hooks"))
				if err != nil {
					return err
				}
				for _, hook := range allHooks {
					if allowed && hooks.AppliesToStagedTag(hook, stagedItemTag) {
						err = hooks.ApplyHookToStagedTag(hook, stagedItemTag)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
