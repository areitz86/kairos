package mos_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/spectrocloud/peg/matcher"
)

var _ = Describe("kairos bundles test", Label("bundles-test"), func() {
	var vm VM
	BeforeEach(func() {
		if os.Getenv("CLOUD_INIT") == "" || !filepath.IsAbs(os.Getenv("CLOUD_INIT")) {
			Fail("CLOUD_INIT must be set and must be pointing to a file as an absolute path")
		}
		_, vm = startVM()
		vm.EventuallyConnects(1200)
	})

	AfterEach(func() {
		vm.Destroy(func(vm VM) {
			gatherLogs(vm)
		})
	})

	Context("reboots and passes functional tests", func() {
		BeforeEach(func() {
			expectDefaultService(vm)
			expectStartedInstallation(vm)
			expectRebootedToActive(vm)
		})

		It("has grubenv file", func() {
			By("checking after-install hook triggered")

			Eventually(func() string {
				out, _ := vm.Sudo("cat /oem/grubenv")
				return out
			}, 40*time.Minute, 1*time.Second).Should(
				Or(
					ContainSubstring("foobarzz"),
				))
		})

		It("has custom cmdline", func() {
			By("waiting reboot and checking cmdline is present")
			Eventually(func() string {
				out, _ := vm.Sudo("cat /proc/cmdline")
				return out
			}, 30*time.Minute, 1*time.Second).Should(
				Or(
					ContainSubstring("foobarzz"),
				))
		})

		It("has kubo extension", func() {
			syset, err := vm.Sudo("systemd-sysext")
			ls, _ := vm.Sudo("ls -liah /usr/local/lib/extensions")
			fmt.Println("LS:", ls)
			Expect(err).ToNot(HaveOccurred())
			Expect(syset).To(ContainSubstring("kubo"))

			ipfsV, err := vm.Sudo("ipfs version")
			Expect(err).ToNot(HaveOccurred())

			Expect(ipfsV).To(ContainSubstring("0.15.0"))
		})
	})
})
