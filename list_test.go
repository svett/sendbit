package sendbit_test

import (
	"errors"
	"fmt"
	"os"

	. "github.com/svett/sendbit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List", func() {
	var client *Client
	var cleanUp bool

	BeforeEach(func() {
		cleanUp = true
		user := os.Getenv("SENDGRID_USER")
		Expect(user).ToNot(BeEmpty())
		pass := os.Getenv("SENDGRID_PASS")
		Expect(pass).ToNot(BeEmpty())

		client = NewClient(user, pass)
	})

	AfterEach(func() {
		if !cleanUp {
			return
		}
		lists, err := client.Lists()
		Expect(err).ToNot(HaveOccurred())

		for _, list := range lists {
			Expect(client.DeleteList(list.Name)).To(Succeed())
		}
	})

	It("creates a list", func() {
		name := RandomString(15)
		Expect(client.CreateList(name)).To(Succeed())

		list, err := client.List(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(list).ToNot(BeNil())
		Expect(list.Name).To(Equal(name))
	})

	Context("when Auth is empty or invalid", func() {
		BeforeEach(func() {
			cleanUp = false
			client = &Client{}
		})

		It("fails to create list", func() {
			err := client.CreateList(RandomString(5))
			Expect(err).To(MatchError(errors.New("sendbit: client.CreateList error: " +
				"The client credentails are missing or invalid.")))
		})

		It("fails to get list", func() {
			_, err := client.List(RandomString(5))
			Expect(err).To(MatchError(errors.New("sendbit: client.List error: The " +
				"client credentails are missing or invalid.")))
		})

		It("fails to delete list", func() {
			err := client.DeleteList(RandomString(5))
			Expect(err).To(MatchError(errors.New("sendbit: client.DeleteList error: " +
				"The client credentails are missing or invalid.")))
		})
	})

	Context("when list name is empty", func() {
		It("fails to create list", func() {
			err := client.CreateList("")
			Expect(err).To(MatchError("sendbit: client.CreateList error: The list name cannot be empty."))
		})

		It("fails to delete list", func() {
			err := client.DeleteList("")
			Expect(err).To(MatchError("sendbit: client.DeleteList error: The list name cannot be empty."))
		})

		It("fails to get a specified list", func() {
			list, err := client.List("")
			Expect(list).To(BeNil())
			Expect(err).To(MatchError("sendbit: client.List error: The list name cannot be empty."))
		})
	})

	It("gets all lists", func() {
		nameOne := RandomString(7)
		nameTwo := RandomString(10)
		Expect(client.CreateList(nameOne)).To(Succeed())
		Expect(client.CreateList(nameTwo)).To(Succeed())

		lists, err := client.Lists()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(lists)).To(Equal(2))
		Expect(lists[0].Name).To(Equal(nameOne))
		Expect(lists[1].Name).To(Equal(nameTwo))
	})

	It("deletes a list", func() {
		name := RandomString(10)
		Expect(client.CreateList(name)).To(Succeed())

		list, err := client.List(name)
		Expect(err).ToNot(HaveOccurred())
		Expect(list).ToNot(BeNil())
		Expect(list.Name).To(Equal(name))

		Expect(client.DeleteList(name)).To(Succeed())

		_, err = client.List(name)
		Expect(err).To(MatchError(fmt.Errorf("sendbit: client.List error: the "+
			"title(s) '%s' do not exist", name)))
	})
})
