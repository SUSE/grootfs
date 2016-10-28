package commands // import "code.cloudfoundry.org/grootfs/commands"

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"code.cloudfoundry.org/grootfs/commands/storepath"
	"code.cloudfoundry.org/grootfs/fetcher/local"
	"code.cloudfoundry.org/grootfs/fetcher/remote"
	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/grootfs/image_puller"
	unpackerpkg "code.cloudfoundry.org/grootfs/image_puller/unpacker"
	storepkg "code.cloudfoundry.org/grootfs/store"
	bundlerpkg "code.cloudfoundry.org/grootfs/store/bundler"
	"code.cloudfoundry.org/grootfs/store/cache_driver"
	"code.cloudfoundry.org/grootfs/store/dependency_manager"
	locksmithpkg "code.cloudfoundry.org/grootfs/store/locksmith"
	"code.cloudfoundry.org/grootfs/store/volume_driver"
	"code.cloudfoundry.org/lager"

	"code.cloudfoundry.org/commandrunner/linux_command_runner"
	errorspkg "github.com/pkg/errors"
	"github.com/urfave/cli"
)

var CreateCommand = cli.Command{
	Name:        "create",
	Usage:       "create [options] <image> <id>",
	Description: "Creates a root filesystem for the provided image.",

	Flags: []cli.Flag{
		cli.Int64Flag{
			Name:  "disk-limit-size-bytes",
			Usage: "Inclusive disk limit (i.e: includes all layers in the filesystem)",
		},
		cli.StringSliceFlag{
			Name:  "uid-mapping",
			Usage: "UID mapping for image translation, e.g.: <Namespace UID>:<Host UID>:<Size>",
		},
		cli.StringSliceFlag{
			Name:  "gid-mapping",
			Usage: "GID mapping for image translation, e.g.: <Namespace UID>:<Host UID>:<Size>",
		},
		cli.StringSliceFlag{
			Name:  "insecure-registry",
			Usage: "Whitelist a private registry",
		},
		cli.BoolFlag{
			Name:  "exclude-image-from-quota",
			Usage: "Set disk limit to be exclusive (i.e.: exluding image layers)",
		},
	},

	Action: func(ctx *cli.Context) error {
		logger := ctx.App.Metadata["logger"].(lager.Logger)
		logger = logger.Session("create")

		if ctx.NArg() != 2 {
			logger.Error("parsing-command", errors.New("invalid arguments"), lager.Data{"args": ctx.Args()})
			return cli.NewExitError(fmt.Sprintf("invalid arguments - usage: %s", ctx.Command.Usage), 1)
		}

		storePath, err := storepath.UserBased(ctx.GlobalString("store"))
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("can't determine the store path: %s", err.Error()), 1)
		}

		image := ctx.Args().First()
		id := ctx.Args().Tail()[0]

		diskLimit := ctx.Int64("disk-limit-size-bytes")
		if diskLimit < 0 {
			err := errors.New("invalid argument: disk limit cannot be negative")
			logger.Error("parsing-command", err)
			return cli.NewExitError(err.Error(), 1)
		}

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

		btrfsVolumeDriver := volume_driver.NewBtrfs(ctx.GlobalString("drax-bin"), storePath)
		bundler := bundlerpkg.NewBundler(btrfsVolumeDriver, storePath)

		runner := linux_command_runner.New()
		idMapper := unpackerpkg.NewIDMapper(runner)
		namespacedCmdUnpacker := unpackerpkg.NewNamespacedUnpacker(runner, idMapper)

		trustedRegistries := ctx.StringSlice("insecure-registry")
		dockerSrc := remote.NewDockerSource(trustedRegistries)

		cacheDriver := cache_driver.NewCacheDriver(storePath)
		remoteFetcher := remote.NewRemoteFetcher(dockerSrc, cacheDriver)

		localFetcher := local.NewLocalFetcher()

		locksmith := locksmithpkg.NewFileSystem(storePath)
		imgPuller := image_puller.NewImagePuller(localFetcher, remoteFetcher, namespacedCmdUnpacker, btrfsVolumeDriver)
		dependencyManager := dependency_manager.NewDependencyManager(
			filepath.Join(storePath, storepkg.META_DIR_NAME, "dependencies"),
		)
		creator := groot.IamCreator(bundler, imgPuller, locksmith, dependencyManager)

		bundle, err := creator.Create(logger, groot.CreateSpec{
			ID:                    id,
			Image:                 image,
			DiskLimit:             diskLimit,
			ExcludeImageFromQuota: ctx.Bool("exclude-image-from-quota"),
			UIDMappings:           uidMappings,
			GIDMappings:           gidMappings,
		})
		if err != nil {
			logger.Error("creating", err)

			humanizedError := tryHumanize(err)
			return cli.NewExitError(humanizedError, 1)
		}

		fmt.Println(bundle.Path)
		return nil
	},
}

func parseIDMappings(args []string) ([]groot.IDMappingSpec, error) {
	mappings := []groot.IDMappingSpec{}

	for _, v := range args {
		var mapping groot.IDMappingSpec
		_, err := fmt.Sscanf(v, "%d:%d:%d", &mapping.NamespaceID, &mapping.HostID, &mapping.Size)
		if err != nil {
			return nil, err
		}
		mappings = append(mappings, mapping)
	}

	return mappings, nil
}

func tryHumanize(err error) string {
	switch e := errorspkg.Cause(err).(type) {
	case *url.Error:
		if _, ok := e.Err.(x509.UnknownAuthorityError); ok {
			return "This registry is insecure. To pull images from this registry, please use the --insecure-registry option."
		}
	case remote.ImageNotFoundErr:
		return "Image does not exist or you do not have permissions to see it."
	}

	return err.Error()
}
