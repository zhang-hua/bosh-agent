package directories_test

import (
	. "github.com/cloudfoundry/bosh-agent/internal/github.com/onsi/ginkgo"
	. "github.com/cloudfoundry/bosh-agent/internal/github.com/onsi/gomega"

	"testing"
)

func TestDirectories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Directories Suite")
}
