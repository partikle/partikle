package v1_test

import (
	. "github.com/partikle/partikle/models/api/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	var (
		user *User
		err  error
	)
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
})
