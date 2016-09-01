package commands

import (
	"os"

	limiterpkg "code.cloudfoundry.org/grootfs/store/volume_driver/drax/limiter"
	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry/gunk/command_runner/linux_command_runner"
	"github.com/urfave/cli"
)

var DestroyCommand = cli.Command{
	Name:        "destroy",
	Usage:       "destroy --volume-path <path>",
	Description: "Destroys the qgroup for the given path.",

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "volume-path",
			Usage: "Path to the volume",
		},
	},

	Action: func(ctx *cli.Context) error {
		logger := lager.NewLogger("drax")
		logger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.DEBUG))

		commandRunner := linux_command_runner.New()
		limiter := limiterpkg.NewBtrfsLimiter(commandRunner)
		err := limiter.DestroyQuotaGroup(logger, ctx.String("volume-path"))
		if err != nil {
			logger.Error("destroying-qgroup", err)
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	},
}