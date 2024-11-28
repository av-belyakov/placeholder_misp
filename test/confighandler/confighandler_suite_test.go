package confighandler

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestConfighandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Confighandler Suite")
}
