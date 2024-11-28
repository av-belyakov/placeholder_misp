package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestdatatime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testdatatime Suite")
}
