package virtual_guest_lifecycle_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	datatypes "github.com/maximilien/softlayer-go/data_types"
	softlayer "github.com/maximilien/softlayer-go/softlayer"
	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = Describe("SoftLayer Virtual Guest Lifecycle", func() {
	var (
		err error

		accountService      softlayer.SoftLayer_Account_Service
		virtualGuestService softlayer.SoftLayer_Virtual_Guest_Service
	)

	BeforeEach(func() {
		accountService, err = testhelpers.CreateAccountService()
		Expect(err).ToNot(HaveOccurred())

		virtualGuestService, err = testhelpers.CreateVirtualGuestService()
		Expect(err).ToNot(HaveOccurred())

		testhelpers.TIMEOUT = 35 * time.Minute
		testhelpers.POLLING_INTERVAL = 10 * time.Second
	})

	Context("SoftLayer_Account#<getSshKeys, getVirtualGuests>", func() {
		It("returns an array of SoftLayer_Virtual_Guest objects", func() {
			virtualGuests, err := accountService.GetVirtualGuests()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(virtualGuests)).To(BeNumerically(">=", 0))
		})

		It("returns an array of SoftLayer_Security_Ssh_Keys objects", func() {
			sshKeys, err := accountService.GetSshKeys()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(sshKeys)).To(BeNumerically(">=", 0))
		})
	})

	Context("SoftLayer_SecuritySshKey#CreateObject and SoftLayer_SecuritySshKey#DeleteObject", func() {
		It("creates the ssh key and verify it is present and then deletes it", func() {
			createdSshKey, _ := testhelpers.CreateTestSshKey()
			testhelpers.WaitForCreatedSshKeyToBePresent(createdSshKey.Id)

			sshKeyService, err := testhelpers.CreateSecuritySshKeyService()
			Expect(err).ToNot(HaveOccurred())

			deleted, err := sshKeyService.DeleteObject(createdSshKey.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(deleted).To(BeTrue())

			testhelpers.WaitForDeletedSshKeyToNoLongerBePresent(createdSshKey.Id)
		})
	})

	Context("SoftLayer_VirtualGuest#CreateObject, SoftLayer_VirtualGuest#GetVirtualGuestPrimaryIpAddress, and SoftLayer_VirtualGuest#DeleteObject", func() {
		It("creates the virtual guest instance and waits for it to be active, get it's IP address, and then delete it", func() {
			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			ipAddress := testhelpers.GetVirtualGuestPrimaryIpAddress(virtualGuest.Id)
			Expect(ipAddress).ToNot(Equal(""))

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
		})

		It("creates the virtual guest instance and waits for it to be active, get it's network VLANS, and then delete it", func() {
			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)

			networkVlans, err := virtualGuestService.GetNetworkVlans(virtualGuest.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(networkVlans)).To(BeNumerically(">", 0))

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
		})

		It("creates the virtual guest and waits for it to be active and checks that the host could create a 1MB disk", func() {
			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			virtualGuestService, err := testhelpers.CreateVirtualGuestService()
			Expect(err).ToNot(HaveOccurred())

			available, err := virtualGuestService.CheckHostDiskAvailability(virtualGuest.Id, 1024)
			Expect(err).ToNot(HaveOccurred())
			Expect(available).To(BeTrue())

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
		})
	})

	Context("SoftLayer_VirtualGuest#CreateObject, SoftLayer_VirtualGuest#rebootSoft, wait for reboot to complete, and SoftLayer_VirtualGuest#DeleteObject", func() {
		It("creates the virtual guest instance, wait for active, SOFT reboots it, wait for RUNNING, then delete it", func() {
			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			virtualGuestService, err := testhelpers.CreateVirtualGuestService()
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("----> will attempt to SOFT reboot virtual guest `%d`\n", virtualGuest.Id)
			rebooted, err := virtualGuestService.RebootSoft(virtualGuest.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(rebooted).To(BeTrue())
			fmt.Printf("----> successfully SOFT rebooted virtual guest `%d`\n", virtualGuest.Id)

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
		})
	})

	Context("SoftLayer_VirtualGuest#CreateObject, SoftLayer_VirtualGuest#rebootHard, wait for reboot to complete, and SoftLayer_VirtualGuest#DeleteObject", func() {
		It("creates the virtual guest instance, wait for active, HARD reboots it, wait for RUNNING, then delete it", func() {
			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			virtualGuestService, err := testhelpers.CreateVirtualGuestService()
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("----> will attempt to HARD reboot virtual guest `%d`\n", virtualGuest.Id)
			rebooted, err := virtualGuestService.RebootHard(virtualGuest.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(rebooted).To(BeTrue())
			fmt.Printf("----> successfully HARD rebooted virtual guest `%d`\n", virtualGuest.Id)

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
		})
	})

	Context("SoftLayer_SecuritySshKey#CreateObject and SoftLayer_VirtualGuest#CreateObject", func() {
		It("creates key, creates virtual guest and adds key to list of VG", func() {
			err = testhelpers.FindAndDeleteTestSshKeys()
			Expect(err).ToNot(HaveOccurred())

			createdSshKey, _ := testhelpers.CreateTestSshKey()
			testhelpers.WaitForCreatedSshKeyToBePresent(createdSshKey.Id)

			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{createdSshKey})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
			testhelpers.DeleteSshKey(createdSshKey.Id)
		})
	})

	Context("SoftLayer_VirtualGuest#CreateObject, SoftLayer_VirtualGuest#setTags, and SoftLayer_VirtualGuest#DeleteObject", func() {
		It("creates the virtual guest instance, wait for active, wait for RUNNING, set some tags, verify that tags are added, then delete it", func() {
			virtualGuest := testhelpers.CreateVirtualGuestAndMarkItTest([]datatypes.SoftLayer_Security_Ssh_Key{})

			testhelpers.WaitForVirtualGuestToBeRunning(virtualGuest.Id)
			testhelpers.WaitForVirtualGuestToHaveNoActiveTransactions(virtualGuest.Id)

			virtualGuestService, err := testhelpers.CreateVirtualGuestService()
			Expect(err).ToNot(HaveOccurred())

			fmt.Printf("----> will attempt to set tags to the virtual guest `%d`\n", virtualGuest.Id)
			tags := []string{"tag0", "tag1", "tag2"}
			tagsWasSet, err := virtualGuestService.SetTags(virtualGuest.Id, tags)
			Expect(err).ToNot(HaveOccurred())
			Expect(tagsWasSet).To(BeTrue())

			fmt.Printf("----> verifying that tags were set the tags virtual guest `%d`\n", virtualGuest.Id)
			tagReferences, err := virtualGuestService.GetTagReferences(virtualGuest.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(tagReferences)).To(Equal(3))

			fmt.Printf("----> verify that each tag was set to virtual guest: `%d`\n", virtualGuest.Id)
			found := false
			for _, tag := range tags {
				for _, tagReference := range tagReferences {
					if tag == tagReference.Tag.Name {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue())
				found = false
			}

			fmt.Printf("----> successfully set the tags and verified tags were set in virtual guest `%d`\n", virtualGuest.Id)

			testhelpers.DeleteVirtualGuest(virtualGuest.Id)
		})
	})
})
