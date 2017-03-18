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
	Context("database operations", func() {
		var (
			id  int64
			err error
		)
		BeforeEach(func() {
			id, err = AddUser(user)
		})
		Describe("AddUser", func() {
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
		Describe("GetUserByID", func() {
			Context("user exists", func() {
				It("returns the *User object matching the id", func() {
					getUser, err := GetUserByID(id)
					Expect(err).NotTo(HaveOccurred())
					Expect(getUser.Username).To(Equal(user.Username))
					Expect(getUser.Created.Unix()).To(Equal(user.Created.Unix()))
					Expect(getUser.PasswordHash).To(Equal(user.PasswordHash))
					Expect(getUser.PasswordSalt).To(Equal(user.PasswordSalt))
				})
			})
			Context("invalid user id", func() {
				It("returns an error", func() {
					_, err := GetUserByID(id + 1)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("no row found"))
				})
			})
		})
		FDescribe("GetAllUser", func() {
			var (
				query  map[string]string
				fields []string
				sortBy []string
				order  []string
				offset int64
				limit  int64

				user1, user2, user3 *User
			)
			BeforeEach(func() {
				query = make(map[string]string)

				user1 = createAndAddUser("testuser1", "testpassword")
				user2 = createAndAddUser("testuser2", "testpassword")
				user3 = createAndAddUser("testuser3", "testpassword")
			})
			Context("empty query", func() {
				It("returns all matches", func() {
					assertExpectedUsers(query, fields, sortBy, order, offset, limit,
						user, user1, user2, user3)
				})
			})
			Context("set query[username] = user.Username", func() {
				BeforeEach(func() {
					query["username"] = user.Username
				})
				It("returns a list with only that user", func() {
					assertExpectedUsers(query, fields, sortBy, order, offset, limit,
						user)
				})
			})
		})
	})
})

func createAndAddUser(username, password string) *User {
	user, err := NewUser(username, password)
	testhelpers.Must(err)
	id, err := AddUser(user)
	testhelpers.Must(err)
	user.ID = id
	return user
}

func assertExpectedUsers(query map[string]string, fields []string, sortBy []string, order []string,
	offset int64, limit int64, expUsers ...*User) {
	users, err := GetAllUser(query, fields, sortBy, order, offset, limit)
	Expect(err).NotTo(HaveOccurred())
	var expectedUsers []*User
	for _, user := range expUsers {
		expectedUsers = append(expectedUsers, user)
	}
	Expect(users).To(HaveLen(len(expectedUsers)))
	for i, item := range users {
		user := item.(User)
		expectedUser := expectedUsers[i]
		Expect(user.ID).To(Equal(expectedUser.ID))
		Expect(user.Username).To(Equal(expectedUser.Username))
		Expect(user.PasswordHash).To(Equal(expectedUser.PasswordHash))
		Expect(user.PasswordSalt).To(Equal(expectedUser.PasswordSalt))
		Expect(user.Created.Equal(expectedUser.Created)).To(BeTrue())
	}
}
