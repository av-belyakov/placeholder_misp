package tmpdata

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlaceholderMisp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PlaceholderMisp Suite")
}
