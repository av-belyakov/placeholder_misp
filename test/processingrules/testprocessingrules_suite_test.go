package testprocessingrules_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestprocessingrules(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testprocessingrules Suite")
}
