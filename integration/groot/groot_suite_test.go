package groot_test

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"code.cloudfoundry.org/grootfs/integration/runner"
	"code.cloudfoundry.org/grootfs/testhelpers"
	"code.cloudfoundry.org/lager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	GrootFSBin string
	DraxBin    string
	StorePath  string
	Runner     *runner.Runner
	storeName  string

	CurrentUserID string
)

const btrfsMountPath = "/mnt/btrfs"

func TestGroot(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		grootFSBin, err := gexec.Build("code.cloudfoundry.org/grootfs")
		Expect(err).NotTo(HaveOccurred())

		return []byte(grootFSBin)
	}, func(data []byte) {
		user, err := user.Current()
		Expect(err).NotTo(HaveOccurred())

		CurrentUserID = user.Uid
		GrootFSBin = string(data)
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		gexec.CleanupBuildArtifacts()
	})

	BeforeEach(func() {
		if os.Getuid() == 0 {
			Skip("This suite is only running as groot")
		}

		storeName = fmt.Sprintf("test-store-%d", GinkgoParallelNode())
		StorePath = path.Join(btrfsMountPath, storeName)
		Expect(os.Mkdir(StorePath, 0700)).NotTo(HaveOccurred())

		var err error
		DraxBin, err = gexec.Build("code.cloudfoundry.org/grootfs/store/volume_driver/drax")
		Expect(err).NotTo(HaveOccurred())
		testhelpers.SuidDrax(DraxBin)

		Runner = &runner.Runner{
			GrootFSBin: GrootFSBin,
			StorePath:  StorePath,
			DraxBin:    DraxBin,
			LogLevel:   lager.DEBUG,
			Stderr:     GinkgoWriter,
		}
	})

	AfterEach(func() {
		testhelpers.CleanUpSubvolumes(btrfsMountPath, storeName)
		Expect(os.RemoveAll(StorePath)).To(Succeed())
	})

	RunSpecs(t, "GrootFS Integration Suite - Running as groot")
}
