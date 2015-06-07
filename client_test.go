package sendbit_test

import (
	"errors"

	. "github.com/svett/sendbit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var nilClient *Client

	It("creates a client successfully", func() {
		client, err := NewClient("user", "pass")
		Expect(err).ToNot(HaveOccurred())
		Expect(client).ToNot(Equal(nilClient))
	})

	Context("when username is empty", func() {
		It("fails to create a client", func() {
			client, err := NewClient("", "pass")
			Expect(client).To(Equal(nilClient))
			Expect(err).To(MatchError(
				errors.New("sendbit: Username argument cannot be empty.")))
		})
	})

	Context("when password is empty", func() {
		It("fails to create a client", func() {
			client, err := NewClient("user", "")
			Expect(client).To(Equal(nilClient))
			Expect(err).To(MatchError(
				errors.New("sendbit: Password argument cannot be empty.")))
		})
	})
})
