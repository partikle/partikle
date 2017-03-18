package v1_test

import (
	. "github.com/partikle/partikle/models/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/partikle/partikle/testhelpers"
)

func init() {
}

var _ = Describe("User", func() {
	var (
		user *User
		err  error
	)
	BeforeSuite(func() {
		testhelpers.Must(testhelpers.InitTestDB())
	})
	AfterSuite(func() {
		testhelpers.Must(testhelpers.DestroyTestDB())
	})
	BeforeEach(func() {
		user, err = NewUser("testuser", "testpassword")
	})
	Describe("New", func() {
		It("should return a *User with hashed password", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(user.PasswordSalt).NotTo(BeEmpty())
			Expect(user.PasswordHash).NotTo(Equal("testpassword"))
		})
	})
	Describe("AddUser", func() {
		Context("user has not been added yet", func() {
			It("returns the user id without error", func() {
				id, err := AddUser(user)
				Expect(err).NotTo(HaveOccurred())
				Expect(id).To(Equal(int64(1)))
			})
		})
	})
})
