package testelasticsearch_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestelasticsearch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testelasticsearch Suite")
}
