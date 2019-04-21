package images

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/containerd/containerd/pkg/progress"
	"github.com/godarch/darch/pkg/repository"
	"github.com/urfave/cli"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "list images",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "print only the image refs",
		},
	},
	Action: func(clicontext *cli.Context) error {
		quiet := clicontext.Bool("quiet")

		repo, err := repository.NewSession(repository.DefaultContainerdSocketLocation)
		if err != nil {
			return err
		}
		defer repo.Close()

		imgs, err := repo.GetImages(context.Background())
		if err != nil {
			return err
		}

		if quiet {
			for _, img := range imgs {
				fmt.Println(img.Name + ":" + img.Tag)
			}
			return nil
		}

		tw := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
		fmt.Fprintln(tw, "REPOSITORY\tTAG\tCREATED\tSIZE\t")
		for _, img := range imgs {
			size, err := repo.GetImageSize(context.Background(), fmt.Sprintf("%s:%s", img.Name, img.Tag))
			if err != nil {
				return err
			}

			fmt.Fprintf(tw, "%v\t%v\t%v\t%v\t\n",
				img.Name,
				img.Tag,
				img.CreatedAt.Format("2006-01-02"),
				progress.Bytes(size))
		}

		return tw.Flush()
	},
}
