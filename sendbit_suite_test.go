package sendbit_test

import (
	"math/rand"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/svett/sendbit"

	"testing"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestSendbit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sendbit Suite")
}

// Creates a new client from Environment variables
func NewClientFromEnv() *sendbit.Client {
	user := os.Getenv("SENDGRID_USER")
	Expect(user).ToNot(BeEmpty())
	pass := os.Getenv("SENDGRID_PASS")
	Expect(pass).ToNot(BeEmpty())
	return sendbit.NewClient(user, pass)
}

// Returns a random srting for particular length
func RandomString(length int) string {
	count := len(letters)
	wordmap := make([]rune, length)
	for index := range wordmap {
		rand.Seed(time.Now().Unix())
		wordmap[index] = letters[rand.Intn(count)]
	}
	return string(wordmap)
}
