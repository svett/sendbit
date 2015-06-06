package sendbit_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestSendbit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sendbit Suite")
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
