package natsinteractions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNatsinteractions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Natsinteractions Suite")
}
