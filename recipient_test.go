package sendbit_test

import (
	"errors"

	. "github.com/svett/sendbit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recipient", func() {
	var (
		client       *Client
		nilRecipient *Recipient
	)
	const list = "sendbit.tmp.list"

	BeforeEach(func() {
		var err error
		client, err = NewClientFromEnv()
		Expect(err).ToNot(HaveOccurred())
		Expect(client.CreateList(list)).To(Succeed())
	})

	AfterEach(func() {
		Expect(client.DeleteList(list)).To(Succeed())
	})

	It("is added successufully", func() {
		recipient := &Recipient{
			Name:  "John Smith",
			Email: "j.smith@example.com",
		}

		Expect(client.AddRecipient(list, recipient)).To(Succeed())

		subscriber, err := client.Recipient(list, recipient.Email)
		Expect(err).ToNot(HaveOccurred())
		Expect(subscriber).To(Equal(recipient))
	})

	Context("when duplicated recipient is added", func() {
		It("fails to subscribe the duplicated recipient", func() {
			recipient := &Recipient{
				Name:  "John Smith",
				Email: "j.smith@example.com",
			}

			Expect(client.AddRecipient(list, recipient)).To(Succeed())
			err := client.AddRecipient(list, recipient)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when is nil", func() {
		It("is not added successfully", func() {
			Expect(client.AddRecipient(list, nil)).To(
				MatchError(errors.New("sendbit: client.AddRecipient error: " +
					"The recipeint is nil or has invalid email.")))
		})
	})

	Context("when does not have email address", func() {
		It("is not added successfully", func() {
			recipient := &Recipient{
				Name: "John Smith",
			}
			Expect(client.AddRecipient(list, recipient)).To(
				MatchError(errors.New("sendbit: client.AddRecipient error: " +
					"The recipeint is nil or has invalid email.")))
		})
	})

	Context("when list is empty", func() {
		It("fails to added recipient", func() {
			recipient := &Recipient{
				Name:  "John Smith",
				Email: "j.smith@example.com",
			}
			Expect(client.AddRecipient("", recipient)).To(
				MatchError(errors.New("sendbit: client.AddRecipient error: " +
					"The list is empty.")))
		})

		It("fails to get recipients", func() {
			_, err := client.Recipients("")
			Expect(err).To(MatchError(errors.New("sendbit: client.Recipients " +
				"error: The list is empty.")))
		})

		It("fails to get recipient count", func() {
			_, err := client.RecipientCount("")
			Expect(err).To(MatchError(errors.New("sendbit: client.RecipientCount " +
				"error: The list is empty.")))
		})

		It("fails to delete recipient", func() {
			Expect(client.DeleteRecipient("", "j.smith@example.com")).To(
				MatchError("sendbit: client.DeleteRecipient error: The list is empty."))
		})
	})

	It("list is loaded successfully", func() {
		recipient := &Recipient{
			Name:  "Morgan Freeman",
			Email: "m.freeman@example.com",
		}

		subscriber := &Recipient{
			Name:  "Mike T.",
			Email: "mike.t@example.com",
		}
		Expect(client.AddRecipient(list, recipient)).To(Succeed())
		Expect(client.AddRecipient(list, subscriber)).To(Succeed())

		subscribers, err := client.Recipients(list)
		Expect(err).ToNot(HaveOccurred())
		Expect(len(subscribers)).To(Equal(2))
		Expect(subscribers[0].Email).To(Equal(recipient.Email))
		Expect(subscribers[1].Email).To(Equal(subscriber.Email))
	})

	It("gets list subscriber's count", func() {
		recipient := &Recipient{
			Name:  "J J",
			Email: "j.j@example.com",
		}

		subscriber := &Recipient{
			Name:  "M J",
			Email: "m.j@example.com",
		}
		Expect(client.AddRecipient(list, recipient)).To(Succeed())
		Expect(client.AddRecipient(list, subscriber)).To(Succeed())

		count, err := client.RecipientCount(list)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(uint64(2)))
	})

	It("is deleted successfully", func() {
		recipient := &Recipient{
			Name:  "John Smith",
			Email: "j.smith@example.com",
		}

		Expect(client.AddRecipient(list, recipient)).To(Succeed())
		Expect(client.DeleteRecipient(list, recipient.Email)).To(Succeed())
		recipient, err := client.Recipient(list, recipient.Email)
		Expect(err).ToNot(HaveOccurred())
		Expect(recipient).To(Equal(nilRecipient))
	})

	Context("when non-existing email is deleted", func() {
		It("fails to delete it", func() {
			Expect(client.DeleteRecipient(list, "no.exists@example.com")).To(
				MatchError("sendbit: client.DeleteRecipient error: " +
					"The recipient is not removed."))
		})
	})
})
