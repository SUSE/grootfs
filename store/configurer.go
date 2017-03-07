package store // import "code.cloudfoundry.org/grootfs/store"

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	errorspkg "github.com/pkg/errors"

	"code.cloudfoundry.org/lager"
)

func ConfigureStore(logger lager.Logger, storePath, driver string, ownerUID, ownerGID int) error {
	logger = logger.Session("ensuring-store", lager.Data{"storePath": storePath})
	logger.Debug("starting")
	defer logger.Debug("ending")

	requiredPaths := []string{
		filepath.Join(storePath, ImageDirName),
		filepath.Join(storePath, VolumesDirName),
		filepath.Join(storePath, CacheDirName),
		filepath.Join(storePath, LocksDirName),
		filepath.Join(storePath, MetaDirName),
		filepath.Join(storePath, TempDirName),
		filepath.Join(storePath, MetaDirName, "dependencies"),
	}

	if err := os.Setenv("TMPDIR", filepath.Join(storePath, TempDirName)); err != nil {
		return errorspkg.Wrap(err, "could not set TMPDIR")
	}

	if err := isDirectory(storePath); err != nil {
		return err
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		if err := os.Mkdir(storePath, 0755); err != nil {
			dir, err1 := os.Lstat(storePath)
			if err1 != nil || !dir.IsDir() {
				return errorspkg.Wrapf(err, "making directory `%s`", storePath)
			}
		}

		if err := os.Chown(storePath, ownerUID, ownerGID); err != nil {
			logger.Error("store-ownership-change-failed", err, lager.Data{"target-uid": ownerUID, "target-gid": ownerGID})
			return errorspkg.Wrapf(err, "changing store owner to %d:%d for path %s", ownerUID, ownerGID, storePath)
		}

		if err := os.Chmod(storePath, 0700); err != nil {
			logger.Error("store-permission-change-failed", err)
			return errorspkg.Wrapf(err, "changing store permissions %s", storePath)
		}
	}

	for _, requiredPath := range requiredPaths {
		if err := createDirectory(logger, requiredPath, ownerUID, ownerGID); err != nil {
			return err
		}
	}

	if requiresWhiteout(driver, ownerUID, ownerGID) {
		if err := createWhiteoutDevice(logger, ownerUID, ownerGID, storePath); err != nil {
			return err
		}

		if err := validateWhiteoutDevice(filepath.Join(storePath, WhiteoutDevice)); err != nil {
			logger.Error("validating-whiteout-device-failed", err)
			return err
		}
	}

	return nil
}

func createDirectory(logger lager.Logger, requiredPath string, ownerUID, ownerGID int) error {
	if err := isDirectory(requiredPath); err != nil {
		return err
	}

	if err := os.Mkdir(requiredPath, 0755); err != nil {
		dir, err1 := os.Lstat(requiredPath)
		if err1 != nil || !dir.IsDir() {
			return errorspkg.Wrapf(err, "making directory `%s`", requiredPath)
		}
	}

	if err := os.Chown(requiredPath, ownerUID, ownerGID); err != nil {
		logger.Error("store-ownership-change-failed", err, lager.Data{"target-uid": ownerUID, "target-gid": ownerGID})
		return errorspkg.Wrapf(err, "changing store owner to %d:%d for path %s", ownerUID, ownerGID, requiredPath)
	}
	return nil
}

func isDirectory(requiredPath string) error {
	if info, err := os.Stat(requiredPath); err == nil {
		if !info.IsDir() {
			return errorspkg.Errorf("path `%s` is not a directory", requiredPath)
		}
	}
	return nil
}

func createWhiteoutDevice(logger lager.Logger, ownerUID, ownerGID int, storePath string) error {
	whiteoutDevicePath := filepath.Join(storePath, WhiteoutDevice)
	if _, err := os.Stat(whiteoutDevicePath); os.IsNotExist(err) {
		if err := syscall.Mknod(whiteoutDevicePath, syscall.S_IFCHR, 0); err != nil {
			if err != nil && !os.IsExist(err) {
				logger.Error("creating-whiteout-device-failed", err, lager.Data{"path": whiteoutDevicePath})
				return errors.Wrapf(err, "failed to create whiteout device %s", whiteoutDevicePath)
			}
		}

		if err := os.Chown(whiteoutDevicePath, ownerUID, ownerGID); err != nil {
			logger.Error("whiteout-device-ownership-change-failed", err, lager.Data{"target-uid": ownerUID, "target-gid": ownerGID})
			return errorspkg.Wrapf(err, "changing store owner to %d:%d for path %s", ownerUID, ownerGID, whiteoutDevicePath)
		}
	}
	return nil
}

func validateWhiteoutDevice(path string) error {
	stat, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return err
	}

	statT := stat.Sys().(*syscall.Stat_t)
	if statT.Rdev != 0 || (stat.Mode()&os.ModeCharDevice) != os.ModeCharDevice {
		return errorspkg.Errorf("the whiteout device file is not a valid device %s", path)
	}

	return nil
}

func requiresWhiteout(driver string, uid, gid int) bool {
	return driver == "overlay-xfs" && uid == 0 && gid == 0
}
