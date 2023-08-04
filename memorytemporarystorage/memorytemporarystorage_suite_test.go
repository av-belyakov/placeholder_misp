package memorytemporarystorage_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMemorytemporarystorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Memorytemporarystorage Suite")
}
