package testupdatestdout_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestupdatestdout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testupdatestdout Suite")
}
