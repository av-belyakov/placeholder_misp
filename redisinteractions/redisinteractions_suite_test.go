package redisinteractions_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRedisinteractions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Redisinteractions Suite")
}
