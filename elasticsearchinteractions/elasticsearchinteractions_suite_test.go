package elasticsearchinteractions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestElasticsearchinteractions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Elasticsearchinteractions Suite")
}
