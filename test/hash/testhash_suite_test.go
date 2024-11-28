package testhash_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTesthash(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testhash Suite")
}
