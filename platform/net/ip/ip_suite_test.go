package ip_test

import (
	. "github.com/cloudfoundry/bosh-agent/internal/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/bosh-agent/internal/github.com/onsi/gomega"

	"testing"
)

func TestPlatform(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ip Suite")
}
