package v1_test

import (
	. "github.com/partikle/partikle/models/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/partikle/partikle/testhelpers"
)

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
		testhelpers.Must(testhelpers.RefreshDBState())
		user, err = NewUser("testuser", "testpassword")
	})
	Describe("New", func() {
		Context("provided credentials are valid", func() {
			It("should return a *User with hashed password", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(user.PasswordSalt).NotTo(BeEmpty())
				Expect(user.PasswordHash).NotTo(Equal("testpassword"))
			})
		})
		Context("username is too short", func() {
			BeforeEach(func() {
				_, err = NewUser("test", "testpassword")
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid username: must be at least 8 characters"))
			})
		})
		Context("password is too short", func() {
			BeforeEach(func() {
				_, err = NewUser("testuser", "pass")
			})
			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid password: must be at least 8 characters"))
			})
		})
	})
	Describe("AddUser", func() {
		var (
			id  int64
			err error
		)
		BeforeEach(func() {
			id, err = AddUser(user)
		})
		Context("user has not been added yet", func() {
			It("returns the user id without error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(id).To(Equal(int64(1)))
			})
		})
		Context("a user with the same username has already been added", func() {
			It("returns an error", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(id).To(Equal(int64(1)))
				_, err = AddUser(user)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIQUE constraint failed"))
			})
		})
	})
})
