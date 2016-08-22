package store_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"code.cloudfoundry.org/grootfs/cloner/clonerfakes"
	"code.cloudfoundry.org/grootfs/store"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bundle", func() {
	var (
		logger lager.Logger

		storePath string

		bundler      *store.Bundler
		volumeDriver *clonerfakes.FakeVolumeDriver
	)

	BeforeEach(func() {
		var err error

		storePath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		volumeDriver = new(clonerfakes.FakeVolumeDriver)
		Expect(os.Mkdir(filepath.Join(storePath, "bundles"), 0777)).To(Succeed())
	})

	JustBeforeEach(func() {
		logger = lagertest.NewTestLogger("test-bunlder")
		bundler = store.NewBundler(storePath, volumeDriver)
	})

	AfterEach(func() {
		Expect(os.RemoveAll(storePath)).To(Succeed())
	})

	Describe("MakeBundle", func() {
		It("should return a bundle directory", func() {
			bundle, err := bundler.MakeBundle(logger, "some-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(bundle.Path()).To(BeADirectory())
		})

		It("should keep the bundles in the same bundle directory", func() {
			someBundle, err := bundler.MakeBundle(logger, "some-id")
			Expect(err).NotTo(HaveOccurred())
			anotherBundle, err := bundler.MakeBundle(logger, "another-id")
			Expect(err).NotTo(HaveOccurred())

			Expect(someBundle.Path()).NotTo(BeEmpty())
			Expect(anotherBundle.Path()).NotTo(BeEmpty())

			bundles, err := ioutil.ReadDir(path.Join(storePath, store.BUNDLES_DIR_NAME))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bundles)).To(Equal(2))
		})

		Context("when calling it with two different ids", func() {
			It("should return two different bundle paths", func() {
				bundle, err := bundler.MakeBundle(logger, "some-id")
				Expect(err).NotTo(HaveOccurred())

				anotherBundle, err := bundler.MakeBundle(logger, "another-id")
				Expect(err).NotTo(HaveOccurred())

				Expect(bundle.Path()).NotTo(Equal(anotherBundle.Path()))
			})
		})

		Context("when using the same id twice", func() {
			It("should return an error", func() {
				_, err := bundler.MakeBundle(logger, "some-id")
				Expect(err).NotTo(HaveOccurred())

				_, err = bundler.MakeBundle(logger, "some-id")
				Expect(err).To(MatchError("bundle for id `some-id` already exists"))
			})
		})

		Context("when the store path does not exist", func() {
			BeforeEach(func() {
				storePath = "/non/existing/store"
			})

			It("should return an error", func() {
				_, err := bundler.MakeBundle(logger, "some-id")
				Expect(err).To(MatchError(ContainSubstring("making bundle path")))
			})
		})
	})

	Describe("DeleteBundle", func() {
		var bundlePath string

		BeforeEach(func() {
			bundlePath = path.Join(storePath, store.BUNDLES_DIR_NAME, "some-id")
			Expect(os.MkdirAll(bundlePath, 0755)).To(Succeed())
			Expect(os.MkdirAll(path.Join(bundlePath, "rootfs"), 0755)).To(Succeed())
			Expect(ioutil.WriteFile(path.Join(bundlePath, "foo"), []byte("hello-world"), 0644)).To(Succeed())
		})

		It("uses the bundle rootfs destroyer to delete the rootfs snapshot", func() {
			Expect(bundler.DeleteBundle(logger, "some-id")).To(Succeed())
			Expect(volumeDriver.DestroyCallCount()).To(Equal(1))

			expectedBundle := store.NewBundle(bundlePath)
			_, rootfsPath := volumeDriver.DestroyArgsForCall(0)
			Expect(rootfsPath).To(Equal(expectedBundle.RootFSPath()))
		})

		It("deletes an existing bundle", func() {
			Expect(bundler.DeleteBundle(logger, "some-id")).To(Succeed())
			Expect(bundlePath).NotTo(BeAnExistingFile())
		})

		Context("when the rootfs path doesn't exist", func() {
			It("doesnt use the bundle rootfs destroyer", func() {
				Expect(os.RemoveAll(path.Join(bundlePath, "rootfs"))).To(Succeed())
				Expect(bundler.DeleteBundle(logger, "some-id")).To(Succeed())
				Expect(volumeDriver.DestroyCallCount()).To(Equal(0))
			})
		})

		Context("when the bundle rootfs destroyer fails", func() {
			It("returns an error", func() {
				volumeDriver.DestroyReturns(errors.New("failed"))

				err := bundler.DeleteBundle(logger, "some-id")
				Expect(err).To(MatchError(ContainSubstring("failed")))
			})
		})

		Context("when bundle does not exist", func() {
			It("returns an error", func() {
				err := bundler.DeleteBundle(logger, "cake")
				Expect(err).To(MatchError(ContainSubstring("bundle path not found")))
			})
		})

		Context("when deleting the folder fails", func() {
			BeforeEach(func() {
				Expect(os.Chmod(bundlePath, 0666)).To(Succeed())
			})

			AfterEach(func() {
				// we need to revert permissions because of the outer AfterEach
				Expect(os.Chmod(bundlePath, 0755)).To(Succeed())
			})

			It("returns an error", func() {
				err := bundler.DeleteBundle(logger, "some-id")
				Expect(err).To(MatchError(ContainSubstring("deleting bundle path")))
			})
		})
	})
})
