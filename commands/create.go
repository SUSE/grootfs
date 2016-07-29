package commands

import (
	"errors"
	"fmt"

	clonerpkg "code.cloudfoundry.org/grootfs/cloner"
	graphpkg "code.cloudfoundry.org/grootfs/graph"
	grootpkg "code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry/gunk/command_runner/linux_command_runner"
	"github.com/urfave/cli"
)

var CreateCommand = cli.Command{
	Name:        "create",
	Usage:       "create --image <image> <id>",
	Description: "Creates a root filesystem for the provided image.",

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "image",
			Usage: "Local path or URL to the image",
		},
		cli.StringSliceFlag{
			Name:  "uid-mapping",
			Usage: "UID mapping for image translation, e.g.: <Namespace UID>:<Host UID>:<Size>",
		},
		cli.StringSliceFlag{
			Name:  "gid-mapping",
			Usage: "GID mapping for image translation, e.g.: <Namespace UID>:<Host UID>:<Size>",
		},
	},

	Action: func(ctx *cli.Context) error {
		logger := ctx.App.Metadata["logger"].(lager.Logger)

		graphPath := ctx.GlobalString("graph")
		imagePath := ctx.String("image")
		if ctx.NArg() != 1 {
			logger.Error("parsing-command", errors.New("id was not specified"))
			return cli.NewExitError("id was not specified", 1)
		}
		id := ctx.Args().First()
		uidMappings, err := parseIDMappings(ctx.StringSlice("uid-mapping"))
		if err != nil {
			err = fmt.Errorf("parsing uid-mapping: %s", err)
			logger.Error("parsing-command", err)
			return cli.NewExitError(err.Error(), 1)
		}
		gidMappings, err := parseIDMappings(ctx.StringSlice("gid-mapping"))
		if err != nil {
			err = fmt.Errorf("parsing gid-mapping: %s", err)
			logger.Error("parsing-command", err)
			return cli.NewExitError(err.Error(), 1)
		}

		graph := graphpkg.NewGraph(graphPath)
		runner := linux_command_runner.New()
		cloner := clonerpkg.NewTarCloner(runner, clonerpkg.NewIDMapper(runner))
		groot := grootpkg.IamGroot(graph, cloner)

		bundlePath, err := groot.Create(logger, grootpkg.CreateSpec{
			ID:          id,
			ImagePath:   imagePath,
			UIDMappings: uidMappings,
			GIDMappings: gidMappings,
		})
		if err != nil {
			logger.Error("making-bundle", err)
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Println(bundlePath)
		return nil
	},
}

func parseIDMappings(args []string) ([]grootpkg.IDMappingSpec, error) {
	mappings := []grootpkg.IDMappingSpec{}

	for _, v := range args {
		var mapping grootpkg.IDMappingSpec
		_, err := fmt.Sscanf(v, "%d:%d:%d", &mapping.NamespaceID, &mapping.HostID, &mapping.Size)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}
