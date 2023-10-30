package testsenderzabbix_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTestsenderzabbix(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testsenderzabbix Suite")
}
