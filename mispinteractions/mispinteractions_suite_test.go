package mispinteractions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMispinteractions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mispinteractions Suite")
}
