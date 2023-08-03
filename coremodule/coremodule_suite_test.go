package coremodule_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCoremodule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Coremodule Suite")
}
